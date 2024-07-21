package handler

import "net/http"

// HealthHandler is a simple handler that returns a 200 status code with an empty body
func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
