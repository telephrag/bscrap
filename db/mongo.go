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
	Col *mongo.Collection
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

	collection := client.Database(config.DBName).Collection(config.CollectionName)

	return &MongoInstance{
		Cli: client,
		Col: collection,
	}, nil
}

func (mi *MongoInstance) StoreRelationData(ctx context.Context, rd *binance.RelationData) (*MongoPayload, error) {

	pl := NewMongoPayload(rd)

	ior, err := mi.Col.InsertOne(ctx, pl)
	if err != nil {
		return nil, err
	}

	if ior.InsertedID == nil {
		return nil, errors.New("document was not inserted")
	}

	fod := mi.Col.FindOneAndDelete(ctx, bson.M{"_id": ior.InsertedID})

	var doc MongoPayload
	fod.Decode(&doc)

	return &doc, nil
}
