package binance

import (
	"encoding/json"
)

type CandleStick struct {
	StartTime  int64
	PriceOpen  string
	PriceHigh  string
	PriceLow   string
	PriceClose string
	EndTime    int64
}

func (cs *CandleStick) UnmarshalJSON(rawData []byte) error {
	data := []interface{}{
		&cs.StartTime,
		&cs.PriceOpen,
		&cs.PriceHigh,
		&cs.PriceLow,
		&cs.PriceClose,
		nil,
		&cs.EndTime,
	}

	if err := json.Unmarshal(rawData, &data); err != nil {
		return err
	}

	return nil
}
