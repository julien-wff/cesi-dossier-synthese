package handler

import (
	"crypto/subtle"
	"encoding/json"
	"github.com/julien-wff/cesi-dossier-synthese/internal/utils"
	"net/http"
)

// GetTelemetryHandler returns the telemetry logs as JSON. Needs basic auth.
func GetTelemetryHandler(config *utils.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check basic auth
		user, password, ok := r.BasicAuth()
		userMatch := subtle.ConstantTimeCompare([]byte(user), []byte(config.TelemetryUser))
		passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(config.TelemetryPassword))
		if !ok || userMatch != 1 || passwordMatch != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Telemetry"`)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("{\"error\":\"unauthorized\"}"))
			return
		}

		// Read telemetry logs
		telemetry, err := utils.ReadTelemetry()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("{\"error\":\"error while reading log file\"}"))
			return
		}

		// Compute stats
		stats := utils.ComputeTelemetryStats(telemetry)
		statsJSON, err := json.Marshal(stats)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("{\"error\":\"error while serializing telemetry stats\"}"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"error\":null, \"data\":" + string(statsJSON) + "}"))
	}
}
