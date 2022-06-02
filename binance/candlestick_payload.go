package binance

import (
	"bscrap/config"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CandleStickDataPayload struct {
	ID        interface{}        `json:"_id,omitempty" bson:"_id,omitempty"`
	Symbol    string             `json:"symbol,omitempty" bson:"symbol"`
	Interval  string             `json:"interval,omitempty" bson:"interval"`
	StartTime int64              `json:"startTime,omitempty" bson:"startTime"`
	EndTime   int64              `json:"endTime,omitempty" bson:"endTime"`
	PriceData csdPair            `json:"priceData,omitempty" bson:"priceData"`
	Expire    primitive.DateTime `json:"expire,omitempty" bson:"expire"`
}

type csdPair struct {
	StartTime  []int64  `json:"startTime,omitempty" bson:"startTime"`
	EndTime    []int64  `json:"endTime,omitempty" bson:"endTime"`
	PriceOpen  []string `json:"priceOpen,omitempty" bson:"priceOpen"`
	PriceHigh  []string `json:"priceHigh,omitempty" bson:"priceHigh"`
	PriceLow   []string `json:"priceLow,omitempty" bson:"priceLow"`
	PriceClose []string `json:"priceClose,omitempty" bson:"priceClose"`
}

func newCsdPair(csd *CandleStickData) csdPair {
	l := len(csd.Data)
	pair := csdPair{
		StartTime:  make([]int64, l),
		EndTime:    make([]int64, l),
		PriceOpen:  make([]string, l),
		PriceHigh:  make([]string, l),
		PriceLow:   make([]string, l),
		PriceClose: make([]string, l),
	}

	for i, candle := range csd.Data {
		pair.StartTime[i] = candle.StartTime
		pair.EndTime[i] = candle.EndTime
		pair.PriceOpen[i] = candle.PriceOpen
		pair.PriceHigh[i] = candle.PriceHigh
		pair.PriceLow[i] = candle.PriceLow
		pair.PriceClose[i] = candle.PriceClose
	}

	return pair
}

func NewCandleStickDataPayload(csd *CandleStickData) *CandleStickDataPayload {
	pair := newCsdPair(csd)
	return &CandleStickDataPayload{
		Symbol:    csd.Symbol,
		Interval:  csd.Interval,
		StartTime: pair.StartTime[0],
		EndTime:   pair.EndTime[len(pair.EndTime)-1],
		PriceData: pair,
		Expire:    primitive.NewDateTimeFromTime(time.Now().Add(config.RecordExpirationTime).UTC()),
	}
}

func (pl *CandleStickDataPayload) ToCandleStickData() *CandleStickData {
	csd := &CandleStickData{
		ID:       pl.ID,
		Symbol:   pl.Symbol,
		Interval: pl.Interval,
		Data:     make([]CandleStick, len(pl.PriceData.StartTime)),
	}

	for i := range pl.PriceData.StartTime {
		csd.Data[i].StartTime = pl.PriceData.StartTime[i]
		csd.Data[i].PriceOpen = pl.PriceData.PriceOpen[i]
		csd.Data[i].PriceHigh = pl.PriceData.PriceHigh[i]
		csd.Data[i].PriceLow = pl.PriceData.PriceLow[i]
		csd.Data[i].PriceClose = pl.PriceData.PriceClose[i]
		csd.Data[i].EndTime = pl.PriceData.EndTime[i]
	}

	return csd
}
