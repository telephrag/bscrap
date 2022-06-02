package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (mi *MongoInstance) ReadOneByID(ctx context.Context, colName string, id string) (*mongo.SingleResult, error) {
	col := mi.Col(colName)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := col.FindOne(ctx, bson.M{"_id": objID})
	return res, res.Err()
}
