package bscrap_srv

import (
	"bscrap/binance"
	"bscrap/config"
	"bscrap/db"
	"bscrap/util"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/bson"
)

type Env struct {
	Argv url.Values
	Mi   *db.MongoInstance

	KLDataA *binance.KLineData
	KLDataB *binance.KLineData

	RelData        *binance.RelationData
	RelDataPayload *binance.RelationDataPayload
}

// startTime and endTime are passed in milliseconds (how it's on Binance)
func (env *Env) GetKLD(
	ctx context.Context,
	symbol, interval, limit, startTime, endTime string,
) (*binance.KLineData, error) {

	uri := util.NewURI(config.API_URL, "https").Proceed("klines")
	uri.Symbol(symbol).Interval(interval).Limit(limit).Timeframe(startTime, endTime)
	uriStr, err := uri.String()
	if err != nil {
		return nil, err
	}

	kld, err := env.lookupKLDInDB( // lookup data in db to possibly avoid refering to API
		ctx,
		symbol,
		interval,
		limit,
		startTime,
		endTime,
	)
	if err == nil { // TODO: handle if error is not ErrNoDocument
		kld.FromDB = true
		return kld, nil // return if data is already in db
	}

	resp, err := http.Get(uriStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 { // handle bad request
		var bErr binance.Err
		err = json.Unmarshal(content, &bErr)
		if err != nil {
			return nil, fmt.Errorf("binance %w", err)
		} else {
			return nil, fmt.Errorf("binance %v", bErr)
		}
	}

	kld = &binance.KLineData{}
	if err = json.Unmarshal(content, &kld.Data); err != nil {
		return nil, fmt.Errorf("binance %w", err)
	}
	kld.Symbol = symbol
	kld.Interval = interval

	return kld, nil
}

func (env *Env) lookupKLDInDB(
	ctx context.Context,
	symbol, interval, limit string,
	startTime, endTime string,
) (*binance.KLineData, error) {

	st, et, err := util.DeduceTime(startTime, endTime, interval, limit)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"symbol":   symbol,
		"interval": interval,
		// "interval":  bson.M{"$in": binance.ShorterOrEqualTo(interval)}, // belongs to set of shorter or equal intervals
		"startTime": bson.M{"$exists": true},
		"endTime":   bson.M{"$exists": true},
		"$expr": bson.M{"$and": bson.A{ // startTime and endTime that exist and satisfy expression
			bson.M{"$and": bson.A{ // if given startTime between startTime and endTime of the record
				bson.M{"$lte": bson.A{"$startTime", st}},
				bson.M{"$gte": bson.A{"$endTime", st}},
			}},

			bson.M{"$and": bson.A{ // same for given endTime
				bson.M{"$lte": bson.A{"$startTime", et}},
				bson.M{"$gte": bson.A{"$endTime", et}},
			}},
		}},
	}

	findOneRes := env.Mi.Col(config.BScrapSourceCol).FindOne(ctx, filter)
	if findOneRes.Err() != nil {
		return nil, findOneRes.Err()
	}

	pl := &binance.KLineDataPayload{}
	err = findOneRes.Decode(pl)
	if err != nil {
		return nil, err
	}

	kld, err := pl.ToKLineData().Shrink(st, et)
	if err != nil {
		return nil, err
	}

	return kld, err
}
