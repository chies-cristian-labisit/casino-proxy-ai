package config

import (
	"os"
	"strconv"
)

// DatadogConfig holds APM configuration read from environment variables.
// DD_API_KEY is intentionally absent — it belongs in the Agent sidecar K8s Secret.
type DatadogConfig struct {
	Enabled    bool
	Service    string
	Env        string
	Version    string
	SampleRate float64
}

func loadDatadogConfig() *DatadogConfig {
	enabled := os.Getenv("DD_ENABLED") == "true"

	sampleRate := 1.0
	if s := os.Getenv("DD_TRACE_SAMPLE_RATE"); s != "" {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			sampleRate = v
		}
	}

	return &DatadogConfig{
		Enabled:    enabled,
		Service:    getEnvOrDefault("DD_SERVICE", "ms-casino-go-v2"),
		Env:        getEnvOrDefault("DD_ENV", "development"),
		Version:    getEnvOrDefault("DD_VERSION", "1.0.0"),
		SampleRate: sampleRate,
	}
}
