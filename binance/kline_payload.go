package binance

import (
	"bscrap/config"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KLineDataPayload struct {
	ID        interface{}        `json:"_id,omitempty" bson:"_id,omitempty"`
	Symbol    string             `json:"symbol,omitempty" bson:"symbol"`
	Interval  string             `json:"interval,omitempty" bson:"interval"`
	StartTime int64              `json:"startTime,omitempty" bson:"startTime"`
	EndTime   int64              `json:"endTime,omitempty" bson:"endTime"`
	PriceData kldPair            `json:"priceData,omitempty" bson:"priceData"`
	Expire    primitive.DateTime `json:"expire,omitempty" bson:"expire"`
}

type kldPair struct {
	StartTime  []int64  `json:"startTime,omitempty" bson:"startTime"`
	EndTime    []int64  `json:"endTime,omitempty" bson:"endTime"`
	PriceOpen  []string `json:"priceOpen,omitempty" bson:"priceOpen"`
	PriceHigh  []string `json:"priceHigh,omitempty" bson:"priceHigh"`
	PriceLow   []string `json:"priceLow,omitempty" bson:"priceLow"`
	PriceClose []string `json:"priceClose,omitempty" bson:"priceClose"`
}

func newCsdPair(kld *KLineData) kldPair {
	l := len(kld.Data)
	pair := kldPair{
		StartTime:  make([]int64, l),
		EndTime:    make([]int64, l),
		PriceOpen:  make([]string, l),
		PriceHigh:  make([]string, l),
		PriceLow:   make([]string, l),
		PriceClose: make([]string, l),
	}

	for i, candle := range kld.Data {
		pair.StartTime[i] = candle.StartTime
		pair.EndTime[i] = candle.EndTime
		pair.PriceOpen[i] = candle.PriceOpen
		pair.PriceHigh[i] = candle.PriceHigh
		pair.PriceLow[i] = candle.PriceLow
		pair.PriceClose[i] = candle.PriceClose
	}

	return pair
}

func NewCandleStickDataPayload(kld *KLineData) *KLineDataPayload {
	pair := newCsdPair(kld)
	return &KLineDataPayload{
		Symbol:    kld.Symbol,
		Interval:  kld.Interval,
		StartTime: pair.StartTime[0],
		EndTime:   pair.EndTime[len(pair.EndTime)-1],
		PriceData: pair,
		Expire:    primitive.NewDateTimeFromTime(time.Now().Add(config.RecordLifeTime).UTC()),
	}
}

func (pl *KLineDataPayload) ToKLineData() *KLineData {
	kld := &KLineData{
		ID:       pl.ID,
		Symbol:   pl.Symbol,
		Interval: pl.Interval,
		Data:     make([]KLine, len(pl.PriceData.StartTime)),
	}

	for i := range pl.PriceData.StartTime {
		kld.Data[i].StartTime = pl.PriceData.StartTime[i]
		kld.Data[i].PriceOpen = pl.PriceData.PriceOpen[i]
		kld.Data[i].PriceHigh = pl.PriceData.PriceHigh[i]
		kld.Data[i].PriceLow = pl.PriceData.PriceLow[i]
		kld.Data[i].PriceClose = pl.PriceData.PriceClose[i]
		kld.Data[i].EndTime = pl.PriceData.EndTime[i]
	}

	return kld
}
