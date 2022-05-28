package db

import (
	"bscrap/binance"
	"bscrap/config"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Cli *mongo.Client
}

func InitMongo(uri string) (*MongoInstance, error) {
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return &MongoInstance{
		Cli: client,
	}, nil
}

func (mi MongoInstance) Col(colName string) *mongo.Collection {
	return mi.Cli.Database(config.DBName).Collection(colName)
}

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

	fod := collection.FindOne(ctx, bson.M{"_id": ior.InsertedID})
	// fod := collection.FindOneAndDelete(ctx, bson.M{"_id": ior.InsertedID})

	var doc RelationDataPayload
	fod.Decode(&doc)

	return &doc, nil
}

func (mi MongoInstance) StoreCandleStickData(ctx context.Context, csdA, csdB *binance.CandleStickData) error {

	pl := NewCandleStickDataPayload(csdA, csdB)

	collection := mi.Col(config.SourceDataCol)

	ior, err := collection.InsertOne(ctx, pl)
	if err != nil {
		return err
	}

	if ior.InsertedID == nil {
		return errors.New("document was not inserted")
	}

	// fod := collection.FindOneAndDelete(ctx, bson.M{"_id": ior.InsertedID})
	// if fod.Err() != nil {
	// 	return fod.Err()
	// }

	return nil
}
