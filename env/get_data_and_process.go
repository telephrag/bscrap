package env

import (
	"bscrap/binance"
	"bscrap/util"
	"errors"
	"net/http"
)

func (env *Env) GetDataAndProcess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if len(env.Argv) == 0 {
			util.HttpErrWriter(
				rw,
				errors.New("no uri arguments, at least 3 expected"),
				http.StatusInternalServerError,
			)
			return
		}

		symbolA := env.Argv.Get("symbolA")
		symbolB := env.Argv.Get("symbolB")
		interval := env.Argv.Get("interval")
		limit := env.Argv.Get("limit")
		startTime := env.Argv.Get("startTime")
		endTime := env.Argv.Get("endTime")

		candleStickDataA, err := binance.GetCandleStickData(symbolA, interval, limit, startTime, endTime)
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusInternalServerError)
			return
		}
		processedDataA, err := candleStickDataA.ProcessCandleStickData()
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusBadRequest)
			return
		}

		candleStickDataB, err := binance.GetCandleStickData(symbolB, interval, limit, startTime, endTime)
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusInternalServerError)
			return
		}
		processedDataB, err := candleStickDataB.ProcessCandleStickData()
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusBadRequest)
			return
		}

		rd, err := binance.GetRelation(processedDataA, processedDataB)
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusInternalServerError)
			return
		}

		env.CSDataA = candleStickDataA
		env.CSDataB = candleStickDataB
		env.RData = rd
		next.ServeHTTP(rw, r)
	})
}
