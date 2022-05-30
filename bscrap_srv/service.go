package bscrap_srv

import (
	"bscrap/util"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Run(env *Env) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.With(
		env.CheckMandatoryArgs,
		env.GetDataAndProcess,
		env.Store,
		env.WriteResponse,
	).Get("/", util.Dummy)

	r.With(
		env.Retrieve,
	).Get("/retrieve", util.Dummy)

	return r
}
