package utils

import (
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	Production        bool
	Port              string
	TelemetryUser     string
	TelemetryPassword string
}

func GetAppConfig() *AppConfig {
	env := strings.TrimSpace(strings.ToLower(os.Getenv("APP_ENV")))
	port := os.Getenv("APP_PORT")

	if p, err := strconv.Atoi(port); err != nil || p < 1 || p > 65535 {
		port = "8080"
	}

	return &AppConfig{
		Production:        !strings.HasPrefix(env, "dev"),
		Port:              port,
		TelemetryUser:     os.Getenv("TELEMETRY_USER"),
		TelemetryPassword: os.Getenv("TELEMETRY_PASSWORD"),
	}
}
