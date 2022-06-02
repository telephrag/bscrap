package bscrap_srv

import (
	"bscrap/util"
	"errors"
	"fmt"
	"net/http"
)

func (env *Env) Store(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if err := env.Mi.Cli.Ping(r.Context(), nil); err != nil {
			util.HttpErrWriter(
				rw,
				fmt.Errorf("%w: connection with mongodb does not exist", err),
				http.StatusInternalServerError,
			)
			return
		}

		if env.RData == nil {
			util.HttpErrWriter(
				rw,
				errors.New("RelationData was expected, none received"),
				http.StatusInternalServerError)
			return
		}

		rd, err := env.Mi.StoreRelationData(r.Context(), env.RData)
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusInternalServerError)
			return
		}

		// may overwrite csData
		if !env.CSDataA.FromDB {
			pl, err := env.Mi.StoreCandleStickData(r.Context(), env.CSDataA)
			if err != nil {
				util.HttpErrWriter(rw, err, http.StatusInternalServerError)
				return
			}
			env.CSDataA.ID = pl.ID
		}
		rd.RawDataAID = env.CSDataA.ID

		if !env.CSDataB.FromDB {
			pl, err := env.Mi.StoreCandleStickData(r.Context(), env.CSDataB)
			if err != nil {
				util.HttpErrWriter(rw, err, http.StatusInternalServerError)
				return
			}
			env.CSDataB.ID = pl.ID
		}
		rd.RawDataBID = env.CSDataB.ID

		env.RDataPayload = rd
		next.ServeHTTP(rw, r)
	})
}
