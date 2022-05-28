package services

import (
	"bscrap/env"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func dummy(rw http.ResponseWriter, r *http.Request) {
}

func Handle(env *env.Env) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.With(
		env.CheckMandatoryArgs,
		env.GetDataAndProcess,
		env.Store,
		env.WriteResponse,
	).Get("/", dummy)

	r.Get("/help", nil)

	return r
}
