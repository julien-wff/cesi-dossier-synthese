package handler

import (
	"crypto/subtle"
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

		// Check if telemetry is empty
		if len(*telemetry) == 0 {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("{\"error\":null, \"data\":[]}"))
			return
		}

		// Replace line breaks with commas
		for i, b := range *telemetry {
			if b == '\n' {
				(*telemetry)[i] = ','
			}
		}

		// Strip last comma
		if (*telemetry)[len(*telemetry)-1] == ',' {
			*telemetry = (*telemetry)[:len(*telemetry)-1]
		}

		// Assemble final JSON
		start := []byte("{\"error\":null,\"data\":[")
		end := []byte("]}")
		*telemetry = append(start, append(*telemetry, end...)...)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(*telemetry)
	}
}
