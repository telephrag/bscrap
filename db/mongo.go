package db

import (
	"bscrap/binance"
	"bscrap/config"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstante struct {
	Cli *mongo.Client
	Col *mongo.Collection
}

func InitMongo(uri string) (*MongoInstante, error) {
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

	return &MongoInstante{
		Cli: client,
		Col: collection,
	}, nil
}

func (mi *MongoInstante) StoreRelationData(rd *binance.RelationData) error {

	pl := NewMongoPayload(rd)

	ctx := context.Background()
	ior, err := mi.Col.InsertOne(ctx, pl)
	if err != nil {
		return err
	}

	if ior.InsertedID == nil {
		return errors.New("document was not inserted")
	}

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt

	fod := mi.Col.FindOneAndDelete(ctx, bson.M{"_id": ior.InsertedID})

	var doc MongoPayload
	fod.Decode(&doc)

	fmt.Println(doc)

	return nil
}
