package bscrap_srv

import (
	"bscrap/binance"
	"bscrap/util"
	"errors"
	"net/http"
)

func (env *Env) GetDataAndProcess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if len(env.Argv) == 0 { // do I even need this check if CheckMandatoryArgs exists
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

		klineDataA, err := env.GetKLD(r.Context(), symbolA, interval, limit, startTime, endTime)
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusInternalServerError)
			return
		}
		processedDataA, err := klineDataA.ProcessCandleStickData()
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusBadRequest)
			return
		}

		klineDataB, err := env.GetKLD(r.Context(), symbolB, interval, limit, startTime, endTime)
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusInternalServerError)
			return
		}
		processedDataB, err := klineDataB.ProcessCandleStickData()
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusBadRequest)
			return
		}

		rd, err := binance.GetRelation(processedDataA, processedDataB)
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusInternalServerError)
			return
		}

		env.KLDataA = klineDataA
		env.KLDataB = klineDataB
		env.RelData = rd
		next.ServeHTTP(rw, r)
	})
}
