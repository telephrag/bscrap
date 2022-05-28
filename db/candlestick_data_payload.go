package db

import "bscrap/binance"

type CandleStickDataPayload struct {
	First  csdPair `json:"pair_a,omitempty" bson:"pair_a"`
	Second csdPair `json:"pair_b,omitempty" bson:"pair_b"`
}

type csdPair struct {
	TradeStart []int64  `json:"tradeStart,omitempty" bson:"tradeStart"`
	TradeEnd   []int64  `json:"tradeEnd,omitempty" bson:"tradeEnd"`
	PriceOpen  []string `json:"priceOpen,omitempty" bson:"priceOpen"`
	PriceHigh  []string `json:"priceHigh,omitempty" bson:"priceHigh"`
	PriceLow   []string `json:"priceLow,omitempty" bson:"priceLow"`
	PriceClose []string `json:"priceClose,omitempty" bson:"priceClose"`
}

func newCsdPair(csd *binance.CandleStickData) csdPair {
	l := len(csd.Data)
	pair := csdPair{
		TradeStart: make([]int64, l),
		TradeEnd:   make([]int64, l),
		PriceOpen:  make([]string, l),
		PriceHigh:  make([]string, l),
		PriceLow:   make([]string, l),
		PriceClose: make([]string, l),
	}

	for i, candle := range csd.Data {
		pair.TradeStart[i] = candle.TradeStart
		pair.TradeEnd[i] = candle.TradeEnd
		pair.PriceOpen[i] = candle.PriceOpen
		pair.PriceHigh[i] = candle.PriceHigh
		pair.PriceLow[i] = candle.PriceLow
		pair.PriceClose[i] = candle.PriceClose
	}

	return pair
}

func NewCandleStickDataPayload(a, b *binance.CandleStickData) *CandleStickDataPayload {
	return &CandleStickDataPayload{
		First:  newCsdPair(a),
		Second: newCsdPair(b),
	}
}
