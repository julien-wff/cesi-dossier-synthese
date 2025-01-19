package router

import (
	"github.com/julien-wff/cesi-dossier-synthese/internal/handler"
	"github.com/julien-wff/cesi-dossier-synthese/internal/utils"
	"net/http"
)

// NewRouter returns a new http.Handler that routes requests to the correct handler
func NewRouter(config *utils.AppConfig) http.Handler {
	r := http.NewServeMux()

	// Health check
	r.HandleFunc("GET /health", handler.HealthHandler)

	// Parsing
	r.HandleFunc("POST /api/parse/debug", handler.ParsePdfDebugHandler)
	r.HandleFunc("POST /api/parse", handler.ParsePdfHandler)
	r.HandleFunc("POST /api/share", handler.ParseSharePdfHandler)

	// Static files
	if config.Production {
		r.Handle("/debug/", handler.StaticHtmlHandler("debug"))
		r.Handle("/", handler.StaticFilesHandler())
	}

	return r
}
