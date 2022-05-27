package db

import (
	"bscrap/binance"
)

type MongoPayload struct {
	PairA       pair    `bson:"pair_a"`
	PairB       pair    `bson:"pair_b"`
	Correlation float64 `bson:"correlation"`
	Covariance  float64 `bson:"covariance"`
	TradeStart  int64   `bson:"startTime"`
	TradeEnd    int64   `bson:"endTime"`
	Count       int     `bson:"count"`
}

type pair struct {
	Symbol string  `bson:"symbol"`
	Mean   float64 `bson:"mean"`
	Spread float64 `bson:"spread"`
}

func NewMongoPayload(rd *binance.RelationData) *MongoPayload {
	return &MongoPayload{
		PairA: pair{
			Symbol: rd.First.Symbol,
			Mean:   rd.First.Mean,
			Spread: rd.First.Spread,
		},
		PairB: pair{
			Symbol: rd.Second.Symbol,
			Mean:   rd.Second.Mean,
			Spread: rd.Second.Spread,
		},
		Correlation: rd.Correlation,
		Covariance:  rd.Covariance,
		TradeStart:  rd.First.TradeStart,
		TradeEnd:    rd.First.TradeEnd,
		Count:       rd.First.Count,
	}
}
