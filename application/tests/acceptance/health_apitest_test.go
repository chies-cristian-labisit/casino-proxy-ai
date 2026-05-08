package acceptance_test

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
)

func TestHealth_LivenessReturns200(t *testing.T) {
	apitest.New().
		HandlerFunc(fiberHandler()).
		Get("/liveness").
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestHealth_ReadinessReturns200WhenDatabaseIsReachable(t *testing.T) {
	apitest.New().
		HandlerFunc(fiberHandler()).
		Get("/readiness").
		Expect(t).
		Status(http.StatusOK).
		End()
}
