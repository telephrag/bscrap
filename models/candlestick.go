package models

import (
	"encoding/json"
)

type CandleStickData struct {
	Data []CandleStick
}

type CandleStick struct {
	TradeStart uint64
	PriceOpen  string
	PriceHigh  string
	PriceLow   string
	PriceClose string
	Volume     string
	TradeEnd   uint64
}

func (cs *CandleStick) UnmarshalJSON(rawData []byte) error {
	data := []interface{}{
		&cs.TradeStart,
		&cs.PriceOpen,
		&cs.PriceHigh,
		&cs.PriceLow,
		&cs.PriceClose,
		&cs.Volume,
		&cs.TradeEnd,
	}

	if err := json.Unmarshal(rawData, &data); err != nil {
		return err
	}

	cs.TradeStart /= 1000 // binance stores time in miliseconds hence division by 1000
	cs.TradeEnd /= 1000

	return nil
}
