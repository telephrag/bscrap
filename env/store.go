package env

import (
	"bscrap/util"
	"errors"
	"net/http"
)

func (env *Env) Store(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if err := env.Mi.Cli.Ping(r.Context(), nil); err != nil {
			util.HttpErrWriter(
				rw,
				errors.New("connection with mongodb does not exist"),
				http.StatusInternalServerError,
			)
		}

		if env.RData == nil {
			util.HttpErrWriter(
				rw,
				errors.New("RelationData was expected, none received"),
				http.StatusInternalServerError)
			return
		}

		pl, err := env.Mi.StoreRelationData(r.Context(), env.RData)
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusInternalServerError)
			return
		}

		env.Pl = pl
		next.ServeHTTP(rw, r)
	})
}
