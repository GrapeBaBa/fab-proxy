package pkg

//easyjson:json
type BlockHeight struct {
	Height uint64 `json:"height"`
}

//easyjson:json
type NodeCount struct {
	Count uint `json:"count"`
}

//easyjson:json
type AcceptedTxCount struct {
	Count uint64 `json:"count"`
}

//easyjson:json
type ConfirmedTxCount struct {
	Count uint64 `json:"count"`
}

//easyjson:json
type BlockInfo struct {
	Height     uint64   `json:"height"`
	TxCount    int      `json:"txcount"`
	Hash       string   `json:"hash"`
	PreHash    string   `json:"prehash"`
	CreateTime string   `json:"createtime"`
	TxHashList []string `json:"txhashlist"` // 交易hash列表
}

//easyjson:json
type TxInfo struct {
	TxId string `json:"txid"`
}

//easyjson:json
type InvokeTokenTx struct {
	ContractId    string `json:"contract_id"`
	Method        string `json:"method"`
	Function      string `json:"function"`
	SourceAccount string `json:"source_account"`
	DestAccount   string `json:"dest_account"`
	Amount        int    `json:"amount"`
}

//easyjson:json
type QueryTokenTx struct {
	ContractId string `json:"contract_id"`
	Method     string `json:"method"`
	Function   string `json:"function"`
	Account    string `json:"account"`
}

//easyjson:json
type InvokeTokenContent struct {
	TxContent InvokeTokenTx `json:"txcontent"`
}

//easyjson:json
type QueryTokenContent struct {
	TxContent QueryTokenTx `json:"txcontent"`
}

