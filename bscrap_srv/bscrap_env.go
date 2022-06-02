package bscrap_srv

import (
	"bscrap/binance"
	"bscrap/config"
	"bscrap/db"
	"bscrap/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Env struct {
	Argv url.Values
	Mi   *db.MongoInstance

	CSDataA *binance.CandleStickData
	CSDataB *binance.CandleStickData

	RData        *binance.RelationData
	RDataPayload *binance.RelationDataPayload
}

// startTime and endTime are passed in milliseconds (how it's on Binance)
func (env *Env) GetCSD(
	ctx context.Context,
	symbol, interval, limit, startTime, endTime string,
) (*binance.CandleStickData, error) {

	uri := util.NewURI(config.API_URL, "https").Proceed("klines")
	uri.Symbol(symbol).Interval(interval).Limit(limit).Timeframe(startTime, endTime)
	uriStr, err := uri.String()
	if err != nil {
		return nil, err
	}

	csd, err := env.lookupCSDInDB( // lookup data in db to possibly avoid refering to API
		ctx,
		symbol,
		interval,
		limit,
		startTime,
		endTime,
	)
	if err == nil { // TODO: handle if error is not ErrNoDocument
		csd.FromDB = true
		return csd, nil // return if data is already in db
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

	csd = &binance.CandleStickData{}
	if err = json.Unmarshal(content, &csd.Data); err != nil {
		return nil, fmt.Errorf("binance %w", err)
	}
	csd.Symbol = symbol
	csd.Interval = interval

	return csd, nil
}

func (env *Env) lookupCSDInDB(
	ctx context.Context,
	symbol, interval, limit string,
	startTime, endTime string,
) (*binance.CandleStickData, error) {

	var st, et int64
	var err error

	i := binance.IntervalLengths[interval]
	r, ok := binance.IntervalRemainders[interval]
	if !ok {
		return nil, errors.New("given interval is invalid or not yet supported for db lookup")
	}

	if startTime != "" {
		if st, err = strconv.ParseInt(startTime, 10, 64); err != nil {
			return nil, fmt.Errorf("%w: \"startTime\" must be int64", err)
		}
		x := i - ((st - r) % i)
		st += x
	}

	if endTime != "" {
		if et, err = strconv.ParseInt(endTime, 10, 64); err != nil {
			return nil, fmt.Errorf("%w: \"endTime\" must be int64", err)
		}
		et = et - (et % i) + r + i - 1
	} else { // assume that user wants data up to current moment...
		et = time.Now().UnixMilli() // ...if endTime and limit are not provided
	}

	if limit != "" {
		l, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: \"limit\" must be int64", err)
		}

		temp := st + i*l - 1 // is off by one error possible here? yep, fixed it
		if temp < et || endTime == "" {
			et = temp
		}
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

	findOneRes := env.Mi.Col(config.SourceDataCollection).FindOne(ctx, filter)
	if findOneRes.Err() != nil {
		return nil, findOneRes.Err()
	}

	pl := &binance.CandleStickDataPayload{}
	err = findOneRes.Decode(pl)
	if err != nil {
		return nil, err
	}

	csd, err := pl.ToCandleStickData().Shrink(st, et)
	if err != nil {
		return nil, err
	}

	return csd, err
}
