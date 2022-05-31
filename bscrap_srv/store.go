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

		cs, err := env.Mi.StoreCandleStickData(r.Context(), env.CSDataA, env.CSDataB)
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusInternalServerError)
			return
		}

		rd.RawDataID = cs.ID

		env.RDataPayload = rd
		env.CSDataPayload = cs
		next.ServeHTTP(rw, r)
	})
}
