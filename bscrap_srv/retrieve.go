package bscrap_srv

import (
	"bscrap/config"
	"bscrap/db"
	"bscrap/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (env *Env) Retrieve(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		argv := r.URL.Query()
		if len(argv) > 2 {
			err := errors.New("excessive amount of arguments given. Maximum 6 are allowed")
			util.HttpErrWriter(rw, err, http.StatusBadRequest)
			return
		}

		if err := env.Mi.Cli.Ping(r.Context(), nil); err != nil {
			util.HttpErrWriter(
				rw,
				fmt.Errorf("%w: connection with mongodb does not exist", err),
				http.StatusInternalServerError,
			)
			return
		}

		var data struct {
			Processed *db.RelationDataPayload    `json:"processed"`
			Raw       *db.CandleStickDataPayload `json:"raw"`
		}

		proc := argv.Get("processed")
		if proc != "" {
			res, err := env.Mi.ReadOneByID(r.Context(), config.ResultsCol, proc)
			if err != nil {
				util.HttpErrWriter(
					rw,
					fmt.Errorf("%w: invalid id was given", err),
					http.StatusBadRequest,
				)
				return
			}

			rd := &db.RelationDataPayload{}
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
		}

		raw := argv.Get("raw")
		if raw != "" {
			res, err := env.Mi.ReadOneByID(r.Context(), config.RawDataCol, raw)
			if err != nil {
				util.HttpErrWriter(
					rw,
					fmt.Errorf("%w: invalid id was given", err),
					http.StatusInternalServerError,
				)
				return
			}

			cs := &db.CandleStickDataPayload{}
			err = res.Decode(cs)
			if err != nil {
				util.HttpErrWriter(
					rw,
					fmt.Errorf("%w: invalid id was given", err),
					http.StatusInternalServerError,
				)
				return
			}
			data.Raw = cs
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
