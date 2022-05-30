package bscrap_srv

import (
	"bscrap/util"
	"encoding/json"
	"errors"
	"net/http"
)

func (env *Env) WriteResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if env.RDataPayload == nil {
			util.HttpErrWriter(
				rw,
				errors.New("no payload to write received"),
				http.StatusInternalServerError,
			)
			return
		}

		data, err := json.Marshal(env.RDataPayload)
		if err != nil {
			util.HttpErrWriter(
				rw,
				errors.New("failed to marshal response"),
				http.StatusInternalServerError,
			)
			return
		}

		rw.Header().Set("Content-type", "application/json")
		_, err = rw.Write(data)
		if err != nil {
			util.HttpErrWriter(
				rw,
				errors.New("failed to write response"),
				http.StatusInternalServerError,
			)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
