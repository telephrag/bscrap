package db

import (
	"bscrap/config"
	"context"

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

	return &MongoInstance{
		Cli: client,
	}, nil
}

func (mi MongoInstance) Col(colName string) *mongo.Collection {
	return mi.Cli.Database(config.DBName).Collection(colName)
}
