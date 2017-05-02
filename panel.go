package main

import "net/http"

func apiUpdateSiteFile(w http.ResponseWriter, r *http.Request) string {
	// ereiamjh
	// Would it better to queue upates in a data structure and have root poke at the server?
	return apiResponse{}.String()
}

func apiNewSite(w http.ResponseWriter, r *http.Request) string {
	return apiResponse{}.String()
}
