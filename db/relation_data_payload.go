package db

import (
	"bscrap/binance"
)

type RelationDataPayload struct {
	ID          interface{} `json:"_id,omitempty" bson:"_id,omitempty"`
	RawDataID   interface{} `json:"_raw_data_id,omitempty" bson:"_raw_data_id,omitempty"`
	PairA       rdPair      `json:"pair_a,omitempty" bson:"pair_a"`
	PairB       rdPair      `json:"pair_b,omitempty" bson:"pair_b"`
	Correlation float64     `json:"correlation,omitempty" bson:"correlation"`
	Covariance  float64     `json:"covariance,omitempty" bson:"covariance"`
	StartTime   int64       `json:"startTime,omitempty" bson:"startTime"`
	EndTime     int64       `json:"endTime,omitempty" bson:"endTime"`
	Count       int         `json:"count,omitempty" bson:"count"`
}

type rdPair struct {
	Symbol string  `json:"symbol,omitempty" bson:"symbol"`
	Mean   float64 `json:"mean,omitempty" bson:"mean"`
	Spread float64 `json:"spread,omitempty" bson:"spread"`
}

func NewMongoPayload(rd *binance.RelationData) *RelationDataPayload {
	return &RelationDataPayload{
		PairA: rdPair{
			Symbol: rd.First.Symbol,
			Mean:   rd.First.Mean,
			Spread: rd.First.Spread,
		},
		PairB: rdPair{
			Symbol: rd.Second.Symbol,
			Mean:   rd.Second.Mean,
			Spread: rd.Second.Spread,
		},
		Correlation: rd.Correlation,
		Covariance:  rd.Covariance,
		StartTime:   rd.First.TradeStart,
		EndTime:     rd.First.TradeEnd,
		Count:       rd.First.Count,
	}
}
