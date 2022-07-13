package client

type TransactionsData struct {
	Data []TransactionData    `json:"data"`
	Meta TransactionsMetaData `json:"meta"`
}

type SingleTransactionDataWrapper struct {
	Data TransactionData `json:"data"`
}

type TransactionData struct {
	Type       string                `json:"type"`
	ID         string                `json:"id"`
	Attributes TransactionAttributes `json:"attributes"`
}

type TransactionAttributes struct {
	State    string `json:"state"`
	Data     string `json:"data"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	ChainID  string `json:"evmChainID"`
	GasLimit string `json:"gasLimit"`
	GasPrice string `json:"gasPrice"`
	Hash     string `json:"hash"`
	RawHex   string `json:"rawHex"`
	Nonce    string `json:"nonce"`
	SentAt   string `json:"sentAt"`
}

type TransactionsMetaData struct {
	Count int `json:"count"`
}
