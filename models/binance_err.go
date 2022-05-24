package models

type BinanceErr struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
