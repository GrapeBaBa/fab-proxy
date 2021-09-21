package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"sync"
	"time"

	"github.com/Grapebaba/fab-proxy/pkg/crypto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/protoutil"

	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/comm"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type BroadcastClientWrapper struct {
	Client orderer.AtomicBroadcast_BroadcastClient
	sync.Mutex
}

func (bcw *BroadcastClientWrapper) Lock() {
	bcw.Mutex.Lock()
}

func (bcw *BroadcastClientWrapper) UnLock() {
	bcw.Mutex.Unlock()
}

type DeliverClientWrapper struct {
	Client orderer.AtomicBroadcast_DeliverClient
	sync.Mutex
}

func (dcw *DeliverClientWrapper) Lock() {
	dcw.Mutex.Lock()
}

func (dcw *DeliverClientWrapper) UnLock() {
	dcw.Mutex.Unlock()
}

func CreateGRPCSClient(certPath string) (*comm.GRPCClient, error) {
	caPEM, res := ioutil.ReadFile(certPath)
	if res != nil {
		err := errors.WithMessage(res,
			fmt.Sprintf("unable to load %s cert", certPath))
		return nil, err
	}
	config := comm.ClientConfig{}
	config.KaOpts = comm.DefaultKeepaliveOptions
	config.Timeout = 3 * time.Second
	config.SecOpts = comm.SecureOptions{
		UseTLS:            true,
		RequireClientCert: false,
		ServerRootCAs:     [][]byte{caPEM},
	}

	grpcClient, err := comm.NewGRPCClient(config)
	//to do: unit test for this error, current fails to make case for this
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to %s", "ds")
	}

	return grpcClient, nil
}

func CreateEndorserClient(addr string, certPath string, sn string) (peer.EndorserClient, error) {
	conn, err := DailConnection(addr, certPath, sn)
	if err != nil {
		return nil, err
	}
	return peer.NewEndorserClient(conn), nil
}

func CreateBroadcastClient(addr string, certPath string, sn string) (*BroadcastClientWrapper, error) {
	conn, err := DailConnection(addr, certPath, sn)
	if err != nil {
		return nil, err
	}
	client, err := orderer.NewAtomicBroadcastClient(conn).Broadcast(context.Background())
	if err != nil {
		return nil, err
	}
	return &BroadcastClientWrapper{Client: client}, nil
}

func CreateDeliverFilteredClient(addr string, certPath string, sn string) (peer.Deliver_DeliverFilteredClient, error) {
	conn, err := DailConnection(addr, certPath, sn)
	if err != nil {
		return nil, err
	}
	return peer.NewDeliverClient(conn).DeliverFiltered(context.Background())
}

func CreateDeliverClient(addr string, certPath string, sn string) (*DeliverClientWrapper, error) {
	conn, err := DailConnection(addr, certPath, sn)
	if err != nil {
		return nil, err
	}
	client, err := orderer.NewAtomicBroadcastClient(conn).Deliver(context.Background())
	if err != nil {
		return nil, err
	}
	return &DeliverClientWrapper{Client: client}, nil
}

func DailConnection(addr string, certPath string, sn string) (*grpc.ClientConn, error) {
	gRPCClient, err := CreateGRPCSClient(certPath)
	if err != nil {
		return nil, err
	}
	conn, err := gRPCClient.NewConnection(addr, comm.ServerNameOverride(sn))
	//conn, err := gRPCClient.NewConnection(addr, func(tlsConfig *tls.Config) { tlsConfig.InsecureSkipVerify = true })
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func CreateSignedDeliverNewestEnv(ch string, signer *crypto.Crypto) (*common.Envelope, error) {
	start := &orderer.SeekPosition{
		Type: &orderer.SeekPosition_Newest{
			Newest: &orderer.SeekNewest{},
		},
	}

	stop := &orderer.SeekPosition{
		Type: &orderer.SeekPosition_Specified{
			Specified: &orderer.SeekSpecified{
				Number: math.MaxUint64,
			},
		},
	}

	seekInfo := &orderer.SeekInfo{
		Start:    start,
		Stop:     stop,
		Behavior: orderer.SeekInfo_BLOCK_UNTIL_READY,
	}

	return protoutil.CreateSignedEnvelope(
		common.HeaderType_DELIVER_SEEK_INFO,
		ch,
		signer,
		seekInfo,
		0,
		0,
	)
}
