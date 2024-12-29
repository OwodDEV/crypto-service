package models

type Transaction struct {
	Hash   string
	From   string
	To     string
	Amount string
}

type GetWalletResp struct {
	Balance string `json:"balance"`
}

type GetTransactionResp struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
}
