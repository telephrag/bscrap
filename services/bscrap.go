package services

import (
	"bscrap/db"
	"bscrap/localmw"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func dummy(rw http.ResponseWriter, r *http.Request) {
}

func Handle(mi *db.MongoInstance) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.With(
		localmw.CheckMandatoryArgs,
		localmw.GetDataAndProcess,
		mi.StoreRelationData_MW,
		localmw.WriteResponse,
	).Get("/", dummy)

	r.Get("/help", nil)

	return r
}
