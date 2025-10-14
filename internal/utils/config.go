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
	ProxyHeaders      bool
	RateLimitTokens   uint64
}

func GetAppConfig() *AppConfig {
	env := strings.TrimSpace(strings.ToLower(os.Getenv("APP_ENV")))
	port := os.Getenv("APP_PORT")
	proxyHeaders := strings.TrimSpace(strings.ToLower(os.Getenv("PROXY_HEADERS")))

	if p, err := strconv.Atoi(port); err != nil || p < 1 || p > 65535 {
		port = "8080"
	}

	rateLimitTokens := uint64(20) // Default to 20 tokens / ip / minute
	if t, err := strconv.Atoi(os.Getenv("RATE_LIMIT_TOKENS")); err == nil && t >= 0 {
		rateLimitTokens = uint64(t)
	}

	return &AppConfig{
		Production:        !strings.HasPrefix(env, "dev"),
		Port:              port,
		TelemetryUser:     os.Getenv("TELEMETRY_USER"),
		TelemetryPassword: os.Getenv("TELEMETRY_PASSWORD"),
		ProxyHeaders:      proxyHeaders == "true" || proxyHeaders == "1",
		RateLimitTokens:   rateLimitTokens,
	}
}
