package observability

import (
	"testing"

	"github.com/cometagaming/ms-casino-go-v2/internal/config"
)

func TestInitTracer_DisabledIsNoOp(t *testing.T) {
	cfg := &config.DatadogConfig{Enabled: false}
	stop := InitTracer(cfg)
	if stop == nil {
		t.Fatal("InitTracer should return a non-nil stop function even when disabled")
	}
	// Calling stop on a no-op tracer must not panic.
	stop()
}

func TestInitTracer_EnabledStartsAndStops(t *testing.T) {
	cfg := &config.DatadogConfig{
		Enabled: true,
		Service: "test-service",
		Env:     "test",
		Version: "0.0.1",
	}
	stop := InitTracer(cfg)
	if stop == nil {
		t.Fatal("InitTracer should return a non-nil stop function when enabled")
	}
	// stop must not panic.
	stop()
}
