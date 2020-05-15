package types

import (
	"encoding/json"
	"time"
)

type TransactionResponse struct {
	TotalCount string `json:"total_count"`
	Count      string `json:"count"`
	PageNumber string `json:"page_number"`
	PageTotal  string `json:"page_total"`
	Limit      string `json:"limit"`
	Txs        []Txs  `json:"txs"`
}
type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type Event struct {
	Type       string      `json:"type"`
	Attributes []Attribute `json:"attributes"`
}
type Logs struct {
	MsgIndex int     `json:"msg_index"`
	Log      string  `json:"log"`
	Events   []Event `json:"events"`
}
type Amount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
type TxValue struct {
	Msg        []Msg         `json:"msg"`
	Fee        Fee           `json:"fee"`
	Signatures []TxSignature `json:"signatures"`
	Memo       string        `json:"memo"`
}

type Msg struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}
type Fee struct {
	Amount []interface{} `json:"amount"`
	Gas    string        `json:"gas"`
}
type TxSignature struct {
	PubKey    string `json:"pub_key"`
	Signature string `json:"signature"`
}
type Tx struct {
	Type  string  `json:"type"`
	Value TxValue `json:"value"`
}
type Txs struct {
	Height    string    `json:"height"`
	Txhash    string    `json:"txhash"`
	RawLog    string    `json:"raw_log"`
	Logs      []Logs    `json:"logs"`
	GasWanted string    `json:"gas_wanted"`
	GasUsed   string    `json:"gas_used"`
	Tx        Tx        `json:"tx"`
	Timestamp time.Time `json:"timestamp"`
}
