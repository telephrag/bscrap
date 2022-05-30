package db

import (
	"bscrap/config"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Cli *mongo.Client
}

func ConnectMongo(uri string) (*MongoInstance, error) {
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	mi := MongoInstance{Cli: client}

	index := mongo.IndexModel{
		Keys:    bson.D{{"expire", 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err = mi.Col(config.ResultsCol).Indexes().CreateOne(context.TODO(), index)
	if err != nil {
		return nil, err
	}

	_, err = mi.Col(config.RawDataCol).Indexes().CreateOne(context.TODO(), index)
	if err != nil {
		return nil, err
	}

	return &mi, nil
}

func (mi MongoInstance) Col(colName string) *mongo.Collection {
	return mi.Cli.Database(config.DBName).Collection(colName)
}
