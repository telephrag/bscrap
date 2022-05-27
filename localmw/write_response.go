package localmw

import (
	"bscrap/db"
	"bscrap/util"
	"encoding/json"
	"errors"
	"net/http"
)

func WriteResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		pl, ok := r.Context().Value(util.CtxKey("payload")).(*db.MongoPayload)
		if !ok {
			util.HttpErrWriter(
				rw,
				errors.New("failed to retrieve payload from context"),
				http.StatusInternalServerError,
			)
			return
		}

		data, err := json.Marshal(pl)
		if err != nil {
			util.HttpErrWriter(
				rw,
				errors.New("failed to marshal response"),
				http.StatusInternalServerError,
			)
			return
		}

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
