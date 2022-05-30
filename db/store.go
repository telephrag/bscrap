package db

import (
	"bscrap/binance"
	"bscrap/config"
	"context"
	"errors"
)

func (mi *MongoInstance) StoreRelationData(ctx context.Context, rd *binance.RelationData) (*RelationDataPayload, error) {

	pl := NewMongoPayload(rd)

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

func (mi *MongoInstance) StoreCandleStickData(ctx context.Context, csdA, csdB *binance.CandleStickData) (*CandleStickDataPayload, error) {

	pl := NewCandleStickDataPayload(csdA, csdB)

	collection := mi.Col(config.RawDataCol)

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
