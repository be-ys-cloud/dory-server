package controller

import "net/http"

func HandlePreflight(w http.ResponseWriter, methods string) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Vary", "Origin")
	w.Header().Add("Vary", "Access-Control-Request-Method")
	w.Header().Add("Vary", "Access-Control-Request-Headers")
	w.Header().Add("Access-Control-Expose-Headers", "Authorization, Location")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, Authorization")
	w.Header().Add("Access-Control-Allow-Methods", methods)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept-Content", "application/json")
}

