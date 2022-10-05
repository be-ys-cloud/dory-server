package controller

import (
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"net/http"
)

func SetHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")
		w.Header().Add("Access-Control-Expose-Headers", "Authorization, Location")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, Authorization")
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,OPTIONS,DELETE")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Accept-Content", "application/json")

		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleErrors(err error, w http.ResponseWriter) {
	if err != nil {
		fullError, convertSucceeded := err.(*structures.CustomError)
		if convertSucceeded {
			http.Error(w, fullError.Text, fullError.HttpCode)
		} else {
			http.Error(w, "An unknown (and probably awkward) error occurred.", 500)
		}
	}
}
