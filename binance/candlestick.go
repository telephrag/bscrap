package binance

import (
	"bscrap/config"
	"bscrap/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CandleStickData struct {
	Symbol string
	Data   []candleStick
}

// startTime and endTime are passed in milliseconds (how it's on Binance)
func GetCandleStickData(symbol, interval, limit, startTime, endTime string) (*CandleStickData, error) {

	uri := util.NewURI(config.API_URL, "https").Proceed("klines")
	uri.Symbol(symbol).Interval(interval).Limit(limit).Timeframe(startTime, endTime)
	uriStr, err := uri.String()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(uriStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 { // handle bad request
		var bErr BinanceErr
		err = json.Unmarshal(content, &bErr)
		if err != nil {
			return nil, fmt.Errorf("binance %w", err)
		} else {
			return nil, fmt.Errorf("binance %v", bErr)
		}
	}

	var candleStickData CandleStickData
	if err = json.Unmarshal(content, &candleStickData.Data); err != nil {
		return nil, fmt.Errorf("binance %w", err)
	}
	candleStickData.Symbol = symbol

	return &candleStickData, nil
}

type candleStick struct {
	TradeStart int64
	PriceOpen  string
	PriceHigh  string
	PriceLow   string
	PriceClose string
	TradeEnd   int64
}

func (cs *candleStick) UnmarshalJSON(rawData []byte) error {
	data := []interface{}{
		&cs.TradeStart,
		&cs.PriceOpen,
		&cs.PriceHigh,
		&cs.PriceLow,
		&cs.PriceClose,
		nil,
		&cs.TradeEnd,
	}

	if err := json.Unmarshal(rawData, &data); err != nil {
		return err
	}

	return nil
}
