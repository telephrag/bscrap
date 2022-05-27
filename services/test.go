package services

import (
	"bscrap/binance"
	"bscrap/config"
	"bscrap/db"
	"bscrap/util"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func binanceHandler(rw http.ResponseWriter, r *http.Request) {
	code := http.StatusBadRequest

	args := r.URL.Query()
	if len(args) > 6 {
		err := errors.New("excessive arguments given. Maximum 6 are allowed")
		util.HttpErrWriter(rw, err, code)
		return
	}

	symbolA := args.Get("symbolA") // check if symbols are provided
	symbolB := args.Get("symbolB")
	if symbolA == "" || symbolB == "" {
		err := errors.New("mandatory parameter(s) \"symbolA\", \"symbolB\" are not provided")
		util.HttpErrWriter(rw, err, code)
		return
	}
	if symbolA == symbolB {
		err := errors.New("symbols should be different")
		util.HttpErrWriter(rw, err, code)
		return
	}

	interval := args.Get("interval")
	if interval == "" {
		err := errors.New("mandatory parameter \"interval\" is not provided")
		util.HttpErrWriter(rw, err, code)
		return
	}

	limit := args.Get("limit")

	startTime := args.Get("startTime")
	endTime := args.Get("endTime")

	candleStickDataA, err := binance.GetCandleStickData(symbolA, interval, limit, startTime, endTime)
	if err != nil {
		log.Panic(err)
	}
	processedDataA := candleStickDataA.ProcessCandleStickData()

	candleStickDataB, err := binance.GetCandleStickData(symbolB, interval, limit, startTime, endTime)
	if err != nil {
		log.Panic(err)
	}
	processedDataB := candleStickDataB.ProcessCandleStickData()

	rd, err := binance.GetRelation(processedDataA, processedDataB)
	if err != nil {
		log.Panic(err)
	}

	_mongo, err := db.InitMongo(config.DBUri)
	if err != nil {
		log.Panic(err)
	}

	pl, err := _mongo.StoreRelationData(rd)
	if err != nil {
		log.Panic(err)
	}

	resp, err := json.Marshal(pl)
	if err != nil {
		code = http.StatusInternalServerError
		util.HttpErrWriter(rw, err, code)
		return
	}

	if _, err = rw.Write(resp); err != nil {
		code = http.StatusInternalServerError
		util.HttpErrWriter(rw, err, code)
		return
	}
}

func Handle() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", binanceHandler)

	return r
}
