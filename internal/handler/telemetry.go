package handler

import (
	"crypto/subtle"
	"encoding/json"
	"github.com/julien-wff/cesi-dossier-synthese/internal/utils"
	"net/http"
)

// writeResponse writes the telemetry response (error and data) to the HTTP response writer.
func writeResponse(w http.ResponseWriter, data *utils.TelemetryStats, error string) {
	res := struct {
		Error *string               `json:"error"`
		Data  *utils.TelemetryStats `json:"data"`
	}{
		Error: nil,
		Data:  data,
	}

	if error != "" {
		res.Error = &error
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

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
			writeResponse(w, nil, "Unauthorized access")
			return
		}

		// Read telemetry logs
		telemetry, err := utils.ReadTelemetry()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			writeResponse(w, nil, "Error reading telemetry logs: "+err.Error())
			return
		}

		// Compute stats
		stats := utils.ComputeTelemetryStats(telemetry)
		writeResponse(w, &stats, "")
	}
}
