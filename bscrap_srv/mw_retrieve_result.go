package bscrap_srv

import (
	"bscrap/binance"
	"bscrap/config"
	"bscrap/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (env *Env) Retrieve(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		argv := r.URL.Query()
		if len(argv) > 3 {
			err := errors.New("excessive amount of arguments given. Maximum 3 are allowed")
			util.HttpErrWriter(rw, err, http.StatusBadRequest)
			return
		}

		if err := env.Mi.Cli.Ping(r.Context(), nil); err != nil {
			util.HttpErrWriter(
				rw,
				fmt.Errorf("%w: no connection with database", err),
				http.StatusInternalServerError,
			)
			return
		}

		var data struct {
			Processed *binance.RelationDataPayload `json:"processed,omitempty"`
			RawA      *binance.KLineData           `json:"raw_a,omitempty"`
			RawB      *binance.KLineData           `json:"raw_b,omitempty"`
		}

		rawA := argv.Get("rawA")
		if rawA != "" {
			res, err := env.Mi.ReadOneByID(r.Context(), config.BScrapSourceCol, rawA)
			if err != nil {
				var code int = http.StatusInternalServerError
				switch err {
				case primitive.ErrInvalidHex:
					code = http.StatusBadRequest
				case mongo.ErrNoDocuments:
					code = http.StatusBadRequest
				}

				util.HttpErrWriter(
					rw,
					fmt.Errorf("%w: invalid id was given", err),
					code,
				)
				return
			}

			cs := &binance.KLineDataPayload{}
			err = res.Decode(cs)
			if err != nil {
				var code int = http.StatusInternalServerError
				switch err {
				case primitive.ErrInvalidHex:
					code = http.StatusBadRequest
				case mongo.ErrNoDocuments:
					code = http.StatusBadRequest
				}

				util.HttpErrWriter(
					rw,
					fmt.Errorf("%w: invalid id was given", err),
					code,
				)
				return
			}
			data.RawA = cs.ToKLineData()
		}

		rawB := argv.Get("rawB")
		if rawB != "" {
			res, err := env.Mi.ReadOneByID(r.Context(), config.BScrapSourceCol, rawB)
			if err != nil {
				util.HttpErrWriter(
					rw,
					fmt.Errorf("%w: invalid id was given", err),
					http.StatusInternalServerError,
				)
				return
			}

			cs := &binance.KLineDataPayload{}
			err = res.Decode(cs)
			if err != nil {
				var code int = http.StatusInternalServerError
				switch err {
				case primitive.ErrInvalidHex:
					code = http.StatusBadRequest
				case mongo.ErrNoDocuments:
					code = http.StatusBadRequest
				}

				util.HttpErrWriter(
					rw,
					fmt.Errorf("%w: invalid id was given", err),
					code,
				)
				return
			}

			data.RawB = cs.ToKLineData()
		}

		proc := argv.Get("processed")
		if proc != "" {
			res, err := env.Mi.ReadOneByID(r.Context(), config.BScrapResCol, proc)
			if err != nil {
				util.HttpErrWriter(
					rw,
					fmt.Errorf("%w: invalid id was given", err),
					http.StatusBadRequest,
				)
				return
			}

			rd := &binance.RelationDataPayload{}
			err = res.Decode(rd)
			if err != nil {
				util.HttpErrWriter(
					rw,
					fmt.Errorf("%w: failed to decode processed data", err),
					http.StatusBadRequest,
				)
				return
			}
			data.Processed = rd

			data.RawA, err = data.RawA.ShrinkSelf(data.Processed.StartTime, data.Processed.EndTime)
			if err != nil {
				util.HttpErrWriter(
					rw,
					fmt.Errorf("%w: invalid id was given", err),
					http.StatusInternalServerError,
				)
			}

			data.RawB, err = data.RawB.ShrinkSelf(data.Processed.StartTime, data.Processed.EndTime)
			if err != nil {
				util.HttpErrWriter(
					rw,
					fmt.Errorf("%w: invalid id was given", err),
					http.StatusInternalServerError,
				)
			}
		}

		content, err := json.Marshal(data)
		if err != nil {
			util.HttpErrWriter(
				rw,
				fmt.Errorf("%w: failed to marshal response", err),
				http.StatusInternalServerError,
			)
			return
		}

		rw.Header().Set("Content-type", "application/json")
		rw.Write(content)

		next.ServeHTTP(rw, r)
	})
}
