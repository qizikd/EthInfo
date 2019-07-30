package core

type TxInfo struct {
	Id         int    `json:"id"`
	TxId       string `json:"txid"`
	From       string `json:"from"`
	To         string `json:"to"`
	Token      string `json:"token"`
	Value      int64  `json:"value"`
	GasPrice   int64  `json:"gasPrice"`
	Gas        int64  `json:"gas"`
	GasUse     int64  `json:"gasUse"`
	IsPengding int    `json:"ispengding"`
	TimeStamp  string `json:"time"`
	Status     int    `json:"status"`
	Fee        int64  `json:"fee"`
}
