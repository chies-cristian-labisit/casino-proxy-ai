package config

import (
	"os"
	"testing"
)

func TestDatadogConfig_DisabledByDefault(t *testing.T) {
	os.Unsetenv("DD_ENABLED")
	cfg := loadDatadogConfig()
	if cfg.Enabled {
		t.Error("expected Datadog to be disabled when DD_ENABLED is not set")
	}
}

func TestDatadogConfig_EnabledWhenTrue(t *testing.T) {
	t.Setenv("DD_ENABLED", "true")
	cfg := loadDatadogConfig()
	if !cfg.Enabled {
		t.Error("expected Datadog to be enabled when DD_ENABLED=true")
	}
}

func TestDatadogConfig_DisabledWhenFalse(t *testing.T) {
	t.Setenv("DD_ENABLED", "false")
	cfg := loadDatadogConfig()
	if cfg.Enabled {
		t.Error("expected Datadog to be disabled when DD_ENABLED=false")
	}
}

func TestDatadogConfig_ParsesEnvVars(t *testing.T) {
	t.Setenv("DD_ENABLED", "true")
	t.Setenv("DD_SERVICE", "my-service")
	t.Setenv("DD_ENV", "staging")
	t.Setenv("DD_VERSION", "2.0.0")
	t.Setenv("DD_TRACE_SAMPLE_RATE", "0.5")

	cfg := loadDatadogConfig()

	if cfg.Service != "my-service" {
		t.Errorf("Service: got %q, want %q", cfg.Service, "my-service")
	}
	if cfg.Env != "staging" {
		t.Errorf("Env: got %q, want %q", cfg.Env, "staging")
	}
	if cfg.Version != "2.0.0" {
		t.Errorf("Version: got %q, want %q", cfg.Version, "2.0.0")
	}
	if cfg.SampleRate != 0.5 {
		t.Errorf("SampleRate: got %v, want 0.5", cfg.SampleRate)
	}
}

func TestDatadogConfig_DefaultSampleRate(t *testing.T) {
	os.Unsetenv("DD_TRACE_SAMPLE_RATE")
	cfg := loadDatadogConfig()
	if cfg.SampleRate != 1.0 {
		t.Errorf("SampleRate default: got %v, want 1.0", cfg.SampleRate)
	}
}

func TestDatadogConfig_InvalidSampleRateFallsToDefault(t *testing.T) {
	t.Setenv("DD_TRACE_SAMPLE_RATE", "not-a-number")
	cfg := loadDatadogConfig()
	if cfg.SampleRate != 1.0 {
		t.Errorf("SampleRate with invalid value: got %v, want 1.0", cfg.SampleRate)
	}
}
