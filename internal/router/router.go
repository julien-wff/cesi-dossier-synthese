package router

import (
	"github.com/julien-wff/cesi-dossier-synthese/internal/handler"
	"net/http"
)

// NewRouter returns a new http.Handler that routes requests to the correct handler
func NewRouter() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("GET /health", handler.HealthHandler)

	return r
}
