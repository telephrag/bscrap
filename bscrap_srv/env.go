package bscrap_srv

import (
	"bscrap/binance"
	"bscrap/db"
	"net/url"
)

type Env struct {
	Argv    url.Values
	Mi      *db.MongoInstance
	CSDataA *binance.CandleStickData
	CSDataB *binance.CandleStickData
	RData   *binance.RelationData

	CSDataPayload *db.CandleStickDataPayload
	RDataPayload  *db.RelationDataPayload
}
