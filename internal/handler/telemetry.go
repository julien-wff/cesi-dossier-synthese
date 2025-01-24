package handler

import (
	"github.com/julien-wff/cesi-dossier-synthese/internal/utils"
	"net/http"
)

// GetTelemetryHandler returns the telemetry logs as JSON
func GetTelemetryHandler(w http.ResponseWriter, _ *http.Request) {
	telemetry, err := utils.ReadTelemetry()
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"error\":\"error while reading log file\"}"))
		return
	}

	start := []byte("{\"error\":null,\"data\":[")
	end := []byte("]}")

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

	*telemetry = append(start, append(*telemetry, end...)...)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(*telemetry)
}
