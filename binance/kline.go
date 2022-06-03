package binance

import (
	"encoding/json"
)

type KLine struct {
	StartTime  int64
	PriceOpen  string
	PriceHigh  string
	PriceLow   string
	PriceClose string
	EndTime    int64
}

func (kl *KLine) UnmarshalJSON(rawData []byte) error {
	data := []interface{}{
		&kl.StartTime,
		&kl.PriceOpen,
		&kl.PriceHigh,
		&kl.PriceLow,
		&kl.PriceClose,
		nil,
		&kl.EndTime,
	}

	if err := json.Unmarshal(rawData, &data); err != nil {
		return err
	}

	return nil
}
