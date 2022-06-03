package binance

import (
	"bscrap/config"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RelationDataPayload struct {
	ID          interface{}        `json:"_id,omitempty" bson:"_id,omitempty"`
	RawDataAID  interface{}        `json:"_raw_data_a_id,omitempty" bson:"_raw_data_a_id,omitempty"`
	RawDataBID  interface{}        `json:"_raw_data_b_id,omitempty" bson:"_raw_data_b_id,omitempty"`
	PairA       rdPair             `json:"pair_a,omitempty" bson:"pair_a"`
	PairB       rdPair             `json:"pair_b,omitempty" bson:"pair_b"`
	Correlation float64            `json:"correlation,omitempty" bson:"correlation"`
	Covariance  float64            `json:"covariance,omitempty" bson:"covariance"`
	StartTime   int64              `json:"startTime,omitempty" bson:"startTime"`
	EndTime     int64              `json:"endTime,omitempty" bson:"endTime"`
	Count       int                `json:"count,omitempty" bson:"count"`
	Expire      primitive.DateTime `json:"expire,omitempty" bson:"expire"`
}

type rdPair struct {
	Symbol string  `json:"symbol,omitempty" bson:"symbol"`
	Mean   float64 `json:"mean,omitempty" bson:"mean"`
	Spread float64 `json:"spread,omitempty" bson:"spread"`
	FromDB bool    `json:"fromDB" bson:"fromDB"`
}

func NewRelationDataPayload(rd *RelationData) *RelationDataPayload {
	return &RelationDataPayload{
		PairA: rdPair{
			Symbol: rd.First.Symbol,
			Mean:   rd.First.Mean,
			Spread: rd.First.Spread,
			FromDB: rd.First.FromDB,
		},
		PairB: rdPair{
			Symbol: rd.Second.Symbol,
			Mean:   rd.Second.Mean,
			Spread: rd.Second.Spread,
			FromDB: rd.Second.FromDB,
		},
		Correlation: rd.Correlation,
		Covariance:  rd.Covariance,
		StartTime:   rd.First.TradeStart,
		EndTime:     rd.First.TradeEnd,
		Count:       rd.First.Count,
		Expire:      primitive.NewDateTimeFromTime(time.Now().Add(config.RecordLifeTime).UTC()),
	}
}
