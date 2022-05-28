package env

import (
	"bscrap/util"
	"errors"
	"net/http"
)

func (env *Env) CheckMandatoryArgs(next http.Handler) http.Handler {
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

		env.Argv = argv
		next.ServeHTTP(rw, r)
	})
}
