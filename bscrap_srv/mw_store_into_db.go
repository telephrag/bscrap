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

		if env.RelData == nil {
			util.HttpErrWriter(
				rw,
				errors.New("RelationData was expected, none received"),
				http.StatusInternalServerError)
			return
		}

		rd, err := env.Mi.StoreRelationData(r.Context(), env.RelData)
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusInternalServerError)
			return
		}

		// may overwrite csData
		if !env.KLDataA.FromDB {
			pl, err := env.Mi.StoreKLineData(r.Context(), env.KLDataA)
			if err != nil {
				util.HttpErrWriter(rw, err, http.StatusInternalServerError)
				return
			}
			env.KLDataA.ID = pl.ID
		}
		rd.RawDataAID = env.KLDataA.ID

		if !env.KLDataB.FromDB {
			pl, err := env.Mi.StoreKLineData(r.Context(), env.KLDataB)
			if err != nil {
				util.HttpErrWriter(rw, err, http.StatusInternalServerError)
				return
			}
			env.KLDataB.ID = pl.ID
		}
		rd.RawDataBID = env.KLDataB.ID

		env.RelDataPayload = rd
		next.ServeHTTP(rw, r)
	})
}
