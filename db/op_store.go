package db

import (
	"bscrap/binance"
	"bscrap/config"
	"context"
	"errors"
)

func (mi *MongoInstance) StoreRelationData(ctx context.Context, rd *binance.RelationData) (*binance.RelationDataPayload, error) {

	pl := binance.NewRelationDataPayload(rd)

	collection := mi.Col(config.ResultsCol)

	ior, err := collection.InsertOne(ctx, pl)
	if err != nil {
		return nil, err
	}

	if ior.InsertedID == nil {
		return nil, errors.New("document was not inserted")
	}
	pl.ID = ior.InsertedID

	return pl, nil
}

func (mi *MongoInstance) StoreCandleStickData(
	ctx context.Context,
	csd *binance.CandleStickData,
) (*binance.CandleStickDataPayload, error) {

	pl := binance.NewCandleStickDataPayload(csd)

	collection := mi.Col(config.SourceDataCollection)

	ior, err := collection.InsertOne(ctx, pl)
	if err != nil {
		return nil, err
	}

	if ior.InsertedID == nil {
		return nil, errors.New("document was not inserted")
	}
	pl.ID = ior.InsertedID

	return pl, nil
}
