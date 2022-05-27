package localmw

import (
	"bscrap/binance"
	"bscrap/util"
	"context"
	"errors"
	"net/http"
	"net/url"
)

func GetDataAndProcess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		argv, ok := r.Context().Value(util.CtxKey("argv")).(url.Values)
		if !ok {
			util.HttpErrWriter(
				rw,
				errors.New("arguments were not received from middleware"),
				http.StatusInternalServerError,
			)
			return
		}

		symbolA := argv.Get("symbolA")
		symbolB := argv.Get("symbolB")
		interval := argv.Get("interval")
		limit := argv.Get("limit")
		startTime := argv.Get("startTime")
		endTime := argv.Get("endTime")

		candleStickDataA, err := binance.GetCandleStickData(symbolA, interval, limit, startTime, endTime)
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusInternalServerError)
			return
		}
		processedDataA := candleStickDataA.ProcessCandleStickData()

		candleStickDataB, err := binance.GetCandleStickData(symbolB, interval, limit, startTime, endTime)
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusInternalServerError)
			return
		}
		processedDataB := candleStickDataB.ProcessCandleStickData()

		rd, err := binance.GetRelation(processedDataA, processedDataB)
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(rw, r.WithContext(context.WithValue(
			context.Background(), util.CtxKey("relationData"), rd,
		)))
	})
}
