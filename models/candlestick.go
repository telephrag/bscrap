package models

import (
	"bscrap/config"
	"bscrap/util"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type CandleStickData struct {
	Symbol string
	Data   []candleStick
}

// startTime and endTime are passed in milliseconds (how it's on Binance)
func GetCandleStickData(symbol, interval string, limit int, startTime, endTime uint64) *CandleStickData {
	uri := util.NewURI(config.API_URL, "https").Proceed("klines")
	uri.Symbol(symbol).Interval(interval).Limit(limit).Timeframe(startTime, endTime)
	uriStr, err := uri.String()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get(uriStr)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var candleStickData CandleStickData
	if err = json.Unmarshal(content, &candleStickData.Data); err != nil {
		log.Fatal(err)
	}
	candleStickData.Symbol = symbol

	return &candleStickData
}

type candleStick struct {
	TradeStart uint64
	PriceOpen  string
	PriceHigh  string
	PriceLow   string
	PriceClose string
	Volume     string
	TradeEnd   uint64
}

func (cs *candleStick) UnmarshalJSON(rawData []byte) error {
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
