package db

import (
	"bscrap/binance"
	"bscrap/util"
	"context"
	"errors"
	"net/http"
)

func (mi *MongoInstance) StoreRelationData_MW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		rd, ok := r.Context().Value(util.CtxKey("relationData")).(*binance.RelationData)
		if !ok {
			util.HttpErrWriter(
				rw,
				errors.New("RelationData was expected, none received"),
				http.StatusInternalServerError)
			return
		}

		pl, err := mi.StoreRelationData(r.Context(), rd)
		if err != nil {
			util.HttpErrWriter(rw, err, http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(rw, r.WithContext(
			context.WithValue(context.Background(), util.CtxKey("payload"), pl),
		))
	})
}
