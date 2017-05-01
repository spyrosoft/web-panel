package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func authorize(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if isLoggedIn(w, r) {
			handle(w, r, ps)
		} else {
			serveLoginPage(w)
		}
	}
}

func authorizeApi(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if isLoggedIn(w, r) {
			handle(w, r, ps)
		} else {
			json.NewEncoder(w).Encode(apiResponse{Errors: []string{"Nope."}})
		}
	}
}
