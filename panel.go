package main

import "net/http"

func apiUpdateSiteFile(w http.ResponseWriter, r *http.Request) string {
	return apiResponse{}.String()
}
