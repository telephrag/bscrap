package db

import (
	"bscrap/binance"
	"bscrap/config"
	"context"
	"errors"
)

func (mi *MongoInstance) StoreRelationData(ctx context.Context, rd *binance.RelationData) (*binance.RelationDataPayload, error) {

	pl := binance.NewRelationDataPayload(rd)

	collection := mi.Col(config.BScrapResCol)

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

func (mi *MongoInstance) StoreKLineData(
	ctx context.Context,
	kld *binance.KLineData,
) (*binance.KLineDataPayload, error) {

	pl := binance.NewCandleStickDataPayload(kld)

	collection := mi.Col(config.BScrapSourceCol)

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
