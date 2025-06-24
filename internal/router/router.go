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
	r.HandleFunc("POST /api/parse/debug", utils.RateLimitMiddlewareFunc(handler.ParsePdfDebugHandler))
	r.HandleFunc("POST /api/parse", utils.RateLimitMiddlewareFunc(handler.ParsePdfHandler))
	r.HandleFunc("POST /api/share", utils.RateLimitMiddlewareFunc(handler.ParseSharePdfHandler))

	// Telemetry
	r.HandleFunc("GET /api/telemetry", handler.GetTelemetryHandler(config))

	// Static files
	if config.Production {
		r.Handle("/debug/", handler.StaticHtmlHandler("debug"))
		r.Handle("/telemetry/", handler.StaticHtmlHandler("telemetry"))
		r.Handle("/", handler.StaticFilesHandler())
	}

	return r
}
