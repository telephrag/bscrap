package binance

type Err struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
