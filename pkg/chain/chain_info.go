package chain

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/Grapebaba/fab-proxy/pkg"
	"github.com/Grapebaba/fab-proxy/pkg/client"
	"github.com/Grapebaba/fab-proxy/pkg/config"
	"github.com/Grapebaba/fab-proxy/pkg/crypto"
	"github.com/hyperledger/fabric-protos-go/common"
	lutil "github.com/hyperledger/fabric/common/ledger/util"
	"github.com/hyperledger/fabric/gossip/util"
	"github.com/hyperledger/fabric/protoutil"
)

var serviceInstance *Service
var once sync.Once

type Service struct {
	LatestHeight     uint64
	AcceptedTxTotal  uint64
	ConfirmedTxTotal uint64
	Blocks           sync.Map
	Client           *client.Client
	Crypto           *crypto.Crypto
	Channel          string
	mutex            sync.Mutex
	pendingTxes      sync.Map
	sampleAccountNum int
}

func Init(config *config.Config) {
	once.Do(func() {
		serviceInstance = newChainInfo(config)
	})
}

func GetChainInfo() *Service {
	return serviceInstance
}

func newChainInfo(config *config.Config) *Service {
	clientMgr := client.New(config.Concurrency, config.Orderers)
	cryptoMgr, _ := crypto.New(config.Crypto.MSPID, config.Crypto.PrivKey, config.Crypto.SignCert)
	service := &Service{Client: clientMgr, Crypto: cryptoMgr, Channel: config.Channel, sampleAccountNum: config.AccountNum}
	ctx := context.Background()
	//for _, bClient := range clientMgr.GetBroadcastClients() {
	//	bc := bClient
	//	go func() {
	//		for {
	//			r, err := bc.Client.Recv()
	//			if err != nil {
	//				panic(err.Error())
	//
	//			}
	//			if r.Status == common.Status_SUCCESS {
	//				cb, ok := service.pendingTxes.Load(r.Info)
	//				if !ok {
	//					panic(r.Info)
	//				}
	//				cb.(chan struct{}) <- struct{}{}
	//			} else {
	//				panic(r.Status)
	//			}
	//
	//			select {
	//			case <-ctx.Done():
	//				return
	//			default:
	//
	//			}
	//		}
	//	}()
	//}

	dClient := clientMgr.GetDeliverClient()
	seek, _ := client.CreateSignedDeliverNewestEnv(config.Channel, cryptoMgr)
	_ = dClient.Client.Send(seek)
	//r, err := dClient.Recv()
	//fmt.Println(r)
	//fmt.Println(err)
	go func() {
		for {
			r, err := dClient.Client.Recv()
			if err != nil {
				fmt.Println(err.Error())
			}
			block := r.GetBlock()
			createTimeBytes := block.GetMetadata().GetMetadata()[common.BlockMetadataIndex_ORDERER]
			createTime, _, _ := lutil.DecodeOrderPreservingVarUint64(createTimeBytes)
			//fmt.Println(block.Header.Number)
			//fmt.Println(createTime)
			txCount := len(block.Data.Data)
			txList := make([]string, txCount)
			for i, envBytes := range block.Data.Data {
				envelop, _ := protoutil.UnmarshalEnvelope(envBytes)
				payload, _ := protoutil.UnmarshalPayload(envelop.Payload)
				chhd, _ := protoutil.UnmarshalChannelHeader(payload.Header.GetChannelHeader())
				txList[i] = chhd.TxId
				if common.HeaderType(chhd.Type) == common.HeaderType_CONFIG || common.HeaderType(chhd.Type) == common.HeaderType_ORDERER_TRANSACTION {
					atomic.AddUint64(&service.AcceptedTxTotal, 1)
				}
			}

			blockInfo := pkg.BlockInfo{
				Height:     block.Header.Number + 1,
				TxCount:    len(block.Data.Data),
				PreHash:    hex.EncodeToString(block.Header.PreviousHash),
				CreateTime: strconv.FormatInt(int64(createTime), 10),
				TxHashList: txList,
			}
			service.Blocks.Store(block.Header.Number, blockInfo)
			atomic.StoreUint64(&service.LatestHeight, blockInfo.Height)
			atomic.AddUint64(&service.ConfirmedTxTotal, uint64(txCount))

			select {
			case <-ctx.Done():
				return
			default:

			}

		}
	}()

	service.InitAccounts()
	return service
}

func (ci *Service) GetBlockHeight() pkg.BlockHeight {
	resp := pkg.BlockHeight{Height: atomic.LoadUint64(&ci.LatestHeight)}
	return resp
}

func (ci *Service) GetAcceptedTxCount() pkg.AcceptedTxCount {
	resp := pkg.AcceptedTxCount{Count: atomic.LoadUint64(&ci.AcceptedTxTotal)}
	return resp
}

func (ci *Service) GetConfirmedTxCount() pkg.ConfirmedTxCount {
	return pkg.ConfirmedTxCount{Count: atomic.LoadUint64(&ci.ConfirmedTxTotal)}
}

func (ci *Service) GetBlockInfo(blockNum uint64) pkg.BlockInfo {
	v, exist := ci.Blocks.Load(blockNum)
	if !exist {
		return pkg.BlockInfo{}
	} else {
		return v.(pkg.BlockInfo)
	}
}

func (ci *Service) CreateTx(txKind string) string {
	num := util.RandomInt(ci.sampleAccountNum)
	srcAcc := fmt.Sprintf("account_%d", num)

	if txKind == "query" {
		res := pkg.QueryTokenTx{ContractId: "token", Method: "invoke", Function: "query", Account: srcAcc}
		resCont := pkg.QueryTokenContent{TxContent: res}
		resContStr, _ := resCont.MarshalJSON()
		return string(resContStr)
	}

	var num1 int
	for {
		num1 = util.RandomInt(ci.sampleAccountNum)
		if num1 != num {
			break
		}
	}

	dstAcc := fmt.Sprintf("account_%d", num1)

	res := pkg.InvokeTokenTx{ContractId: "token", Method: "invoke", Function: "transfer", SourceAccount: srcAcc, DestAccount: dstAcc, Amount: 1}
	resCont := pkg.InvokeTokenContent{TxContent: res}
	resContStr, _ := resCont.MarshalJSON()
	return string(resContStr)
}

func (ci *Service) SendTx(txContent string) pkg.TxInfo {
	nonce, _ := protoutil.CreateNonce()
	creator, _ := ci.Crypto.Serialize()
	txid := protoutil.ComputeTxID(nonce, creator)

	txType := common.HeaderType_ENDORSER_TRANSACTION
	chdr := &common.ChannelHeader{
		Type:      int32(txType),
		ChannelId: ci.Channel,
		TxId:      txid,
		Epoch:     uint64(0),
	}

	shdr := &common.SignatureHeader{
		Creator: creator,
		Nonce:   nonce,
	}

	payload := &common.Payload{
		Header: &common.Header{
			ChannelHeader:   protoutil.MarshalOrPanic(chdr),
			SignatureHeader: protoutil.MarshalOrPanic(shdr),
		},
		Data: []byte(txContent),
	}

	payloadBytes, _ := protoutil.GetBytesPayload(payload)

	signature, _ := ci.Crypto.Sign(payloadBytes)

	envelope := &common.Envelope{
		Payload:   payloadBytes,
		Signature: signature,
	}

	//recvChan := make(chan struct{})
	//ci.pendingTxes.Store(txid, recvChan)

	sender := ci.Client.GetBroadcastClient()
	sender.Lock()
	err := sender.Client.Send(envelope)
	r, err := sender.Client.Recv()
	sender.Unlock()

	if err != nil {
		panic(err)
	}

	if r.Status != common.Status_SUCCESS {
		panic(r.Info)
	}

	//<-recvChan
	atomic.AddUint64(&ci.AcceptedTxTotal, 1)
	return pkg.TxInfo{
		TxId: txid,
	}
}

func (ci *Service) SendTxQuery(txContent string) string {
	nonce, _ := protoutil.CreateNonce()
	creator, _ := ci.Crypto.Serialize()
	txid := protoutil.ComputeTxID(nonce, creator)

	txType := common.HeaderType_ENDORSER_TRANSACTION
	chdr := &common.ChannelHeader{
		Type:      int32(txType),
		ChannelId: ci.Channel,
		TxId:      txid,
		Epoch:     uint64(0),
		Extension: []byte("read"),
	}

	shdr := &common.SignatureHeader{
		Creator: creator,
		Nonce:   nonce,
	}

	payload := &common.Payload{
		Header: &common.Header{
			ChannelHeader:   protoutil.MarshalOrPanic(chdr),
			SignatureHeader: protoutil.MarshalOrPanic(shdr),
		},
		Data: []byte(txContent),
	}

	payloadBytes, _ := protoutil.GetBytesPayload(payload)

	signature, _ := ci.Crypto.Sign(payloadBytes)

	envelope := &common.Envelope{
		Payload:   payloadBytes,
		Signature: signature,
	}

	sender := ci.Client.GetRBroadcastClient()
	sender.Lock()
	err := sender.Client.Send(envelope)
	r, err := sender.Client.Recv()
	sender.Unlock()

	if err != nil {
		panic(err)
	}

	return r.Info
}

func (ci *Service) InitAccounts() pkg.TxInfo {
	nonce, _ := protoutil.CreateNonce()
	creator, _ := ci.Crypto.Serialize()
	txid := protoutil.ComputeTxID(nonce, creator)

	txType := common.HeaderType_ENDORSER_TRANSACTION
	chdr := &common.ChannelHeader{
		Type:      int32(txType),
		ChannelId: ci.Channel,
		TxId:      txid,
		Epoch:     uint64(0),
	}

	shdr := &common.SignatureHeader{
		Creator: creator,
		Nonce:   nonce,
	}

	payload := &common.Payload{
		Header: &common.Header{
			ChannelHeader:   protoutil.MarshalOrPanic(chdr),
			SignatureHeader: protoutil.MarshalOrPanic(shdr),
		},
		Data: []byte(fmt.Sprintf("token,invoke,open,%d", ci.sampleAccountNum)),
	}

	payloadBytes, _ := protoutil.GetBytesPayload(payload)

	signature, _ := ci.Crypto.Sign(payloadBytes)

	envelope := &common.Envelope{
		Payload:   payloadBytes,
		Signature: signature,
	}

	bc := ci.Client.GetInitBroadcastClient()
	bc.Lock()
	err := bc.Client.Send(envelope)
	recv, err := bc.Client.Recv()
	bc.Unlock()

	if err != nil {
		panic(err)
	}

	if recv.Status != common.Status_SUCCESS {
		fmt.Println(recv.Info)
	}
	atomic.AddUint64(&ci.AcceptedTxTotal, 1)

	return pkg.TxInfo{
		TxId: txid,
	}
}

func (ci *Service) CreateTx1(txKind string) string {
	num := util.RandomInt(ci.sampleAccountNum)
	srcAcc := fmt.Sprintf("account_%d", num)

	if txKind == "query" {
		nonce, _ := protoutil.CreateNonce()
		creator, _ := ci.Crypto.Serialize()
		txid := protoutil.ComputeTxID(nonce, creator)

		txType := common.HeaderType_ENDORSER_TRANSACTION
		chdr := &common.ChannelHeader{
			Type:      int32(txType),
			ChannelId: ci.Channel,
			TxId:      txid,
			Epoch:     uint64(0),
			Extension: []byte("read"),
		}

		shdr := &common.SignatureHeader{
			Creator: creator,
			Nonce:   nonce,
		}

		payload := &common.Payload{
			Header: &common.Header{
				ChannelHeader:   protoutil.MarshalOrPanic(chdr),
				SignatureHeader: protoutil.MarshalOrPanic(shdr),
			},
			Data: []byte(fmt.Sprintf("token,invoke,query,%s", srcAcc)),
		}

		payloadBytes, _ := protoutil.GetBytesPayload(payload)

		signature, _ := ci.Crypto.Sign(payloadBytes)

		envelope := &common.Envelope{
			Payload:   payloadBytes,
			Signature: signature,
		}

		envBytes := protoutil.MarshalOrPanic(envelope)

		res := hex.EncodeToString(envBytes)
		resCont := pkg.TokenTx{TxContent: res}
		resContStr, _ := resCont.MarshalJSON()
		return string(resContStr)
	}

	var num1 int
	for {
		num1 = util.RandomInt(ci.sampleAccountNum)
		if num1 != num {
			break
		}
	}

	dstAcc := fmt.Sprintf("account_%d", num1)

	nonce, _ := protoutil.CreateNonce()
	creator, _ := ci.Crypto.Serialize()
	txid := protoutil.ComputeTxID(nonce, creator)

	txType := common.HeaderType_ENDORSER_TRANSACTION
	chdr := &common.ChannelHeader{
		Type:      int32(txType),
		ChannelId: ci.Channel,
		TxId:      txid,
		Epoch:     uint64(0),
	}

	shdr := &common.SignatureHeader{
		Creator: creator,
		Nonce:   nonce,
	}

	payload := &common.Payload{
		Header: &common.Header{
			ChannelHeader:   protoutil.MarshalOrPanic(chdr),
			SignatureHeader: protoutil.MarshalOrPanic(shdr),
		},
		Data: []byte(fmt.Sprintf("token,invoke,transfer,%s,%s,%d", srcAcc, dstAcc, 1)),
	}

	payloadBytes, _ := protoutil.GetBytesPayload(payload)

	signature, _ := ci.Crypto.Sign(payloadBytes)

	envelope := &common.Envelope{
		Payload:   payloadBytes,
		Signature: signature,
	}

	envBytes := protoutil.MarshalOrPanic(envelope)

	res := hex.EncodeToString(envBytes)
	resCont := pkg.TokenTx{TxContent: res}
	resContStr, _ := resCont.MarshalJSON()
	return string(resContStr)
}

func (ci *Service) SendTx1(txContent string) pkg.TxInfo {
	envBytes, _ := hex.DecodeString(txContent)
	envelope, _ := protoutil.UnmarshalEnvelope(envBytes)
	//recvChan := make(chan struct{})
	//ci.pendingTxes.Store(txid, recvChan)

	sender := ci.Client.GetBroadcastClient()
	sender.Lock()
	err := sender.Client.Send(envelope)
	r, err := sender.Client.Recv()
	sender.Unlock()

	if err != nil {
		panic(err)
	}

	if r.Status != common.Status_SUCCESS {
		panic(r.Info)
	}

	//<-recvChan
	atomic.AddUint64(&ci.AcceptedTxTotal, 1)
	return pkg.TxInfo{
		TxId: r.Info,
	}
}

func (ci *Service) SendTxQuery1(txContent string) string {
	envBytes, _ := hex.DecodeString(txContent)
	envelope, _ := protoutil.UnmarshalEnvelope(envBytes)

	sender := ci.Client.GetRBroadcastClient()
	sender.Lock()
	err := sender.Client.Send(envelope)
	r, err := sender.Client.Recv()
	sender.Unlock()

	if err != nil {
		panic(err)
	}

	return r.Info
}