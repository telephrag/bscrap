package localmw

import (
	"bscrap/util"
	"context"
	"errors"
	"net/http"
)

func CheckMandatoryArgs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		argv := r.URL.Query()
		if len(argv) > 6 {
			err := errors.New("excessive amount of arguments given. Maximum 6 are allowed")
			util.HttpErrWriter(rw, err, http.StatusBadRequest)
			return
		}

		symbolA := argv.Get("symbolA")
		symbolB := argv.Get("symbolB")
		if symbolA == "" || symbolB == "" {
			err := errors.New("mandatory parameters \"symbolA\", \"symbolB\" both must be provided")
			util.HttpErrWriter(rw, err, http.StatusBadRequest)
			return
		}
		if symbolA == symbolB {
			err := errors.New("symbols should be different")
			util.HttpErrWriter(rw, err, http.StatusBadRequest)
			return
		}

		interval := argv.Get("interval")
		if interval == "" {
			err := errors.New("mandatory parameter \"interval\" must be provided")
			util.HttpErrWriter(rw, err, http.StatusBadRequest)
			return
		}

		next.ServeHTTP(rw, r.WithContext(
			context.WithValue(r.Context(), util.CtxKey("argv"), argv),
		))
	})
}
