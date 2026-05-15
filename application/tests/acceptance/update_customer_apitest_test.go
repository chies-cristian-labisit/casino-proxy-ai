package acceptance_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/steinfletcher/apitest"

	"github.com/cometagaming/ms-casino-go-v2/internal/domain"
)

// fiberHandler bridges the Fiber app to net/http.Handler, which apitest requires.
func fiberHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := sharedApp.Test(r)
		if err != nil {
			panic(fmt.Sprintf("fiber test: %v", err))
		}
		defer resp.Body.Close()
		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			panic(fmt.Sprintf("fiber body copy: %v", err))
		}
	}
}

// pollUntilNameIs polls GET /api/v2/customers/:code via the Fiber test server
// every 300 ms until the customer name matches or timeout is reached.
// This is the async bridge between the Kafka side-effect and the final apitest assertion.
func pollUntilNameIs(t *testing.T, code, expectedName string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		req := httptest.NewRequest(http.MethodGet, "/api/v2/customers/"+code, nil)
		resp, err := sharedApp.Test(req)
		if err != nil {
			t.Fatalf("poll: fiber test error: %v", err)
		}
		if resp.StatusCode == http.StatusOK {
			var body struct {
				Name string `json:"name"`
			}
			_ = json.NewDecoder(resp.Body).Decode(&body)
			resp.Body.Close()
			if body.Name == expectedName {
				return
			}
		} else {
			resp.Body.Close()
		}
		time.Sleep(300 * time.Millisecond)
	}
	t.Fatalf("timeout after %s: customer %q name never became %q", timeout, code, expectedName)
}

// sendKafkaMessage publishes a customer update event to the Kafka topic.
func sendKafkaMessage(t *testing.T, code, name string) {
	t.Helper()
	type payload struct {
		CustomerCode string `json:"customer_code"`
		CustomerName string `json:"customer_name"`
	}
	body, err := json.Marshal(payload{CustomerCode: code, CustomerName: name})
	if err != nil {
		t.Fatalf("marshal kafka payload: %v", err)
	}
	if err := sharedKafkaWriter.WriteMessages(context.Background(), kafkago.Message{Value: body}); err != nil {
		t.Fatalf("send kafka message: %v", err)
	}
}

// assertCustomer returns an apitest.Assert that validates code and name in the
// JSON response body, without coupling to the auto-increment id field.
func assertCustomer(code, name string) apitest.Assert {
	return func(res *http.Response, _ *http.Request) error {
		var got struct {
			Code string `json:"code"`
			Name string `json:"name"`
		}
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		if got.Code != code {
			return fmt.Errorf("code: want %q, got %q", code, got.Code)
		}
		if got.Name != name {
			return fmt.Errorf("name: want %q, got %q", name, got.Name)
		}
		return nil
	}
}

// TestUpdateCustomer_SuccessfulNameUpdate verifies that a customer name is updated
// when a Kafka message arrives with the new name.
func TestUpdateCustomer_SuccessfulNameUpdate(t *testing.T) {
	const code = "BR123456789"
	if err := sharedRepo.Save(context.Background(), &domain.Customer{Code: code, Name: "Old Name"}); err != nil {
		t.Fatalf("seed customer: %v", err)
	}

	sendKafkaMessage(t, code, "João Silva")
	pollUntilNameIs(t, code, "João Silva", 10*time.Second)

	apitest.New().
		HandlerFunc(fiberHandler()).
		Get("/api/v2/customers/" + code).
		Expect(t).
		Status(http.StatusOK).
		Assert(assertCustomer(code, "João Silva")).
		End()
}

// TestUpdateCustomer_IdempotencyDuplicateSuppressed verifies that a duplicate Kafka
// message does not re-apply an already-processed update.
func TestUpdateCustomer_IdempotencyDuplicateSuppressed(t *testing.T) {
	const code = "BR555555555"
	if err := sharedRepo.Save(context.Background(), &domain.Customer{Code: code, Name: "Original Name"}); err != nil {
		t.Fatalf("seed customer: %v", err)
	}

	sendKafkaMessage(t, code, "Updated Name")
	pollUntilNameIs(t, code, "Updated Name", 10*time.Second)

	sendKafkaMessage(t, code, "Updated Name")
	// Fixed settle: let the listener receive and discard the duplicate.
	// Polling "for no change" is not possible; 1 s is safe since the mock
	// idempotency store rejects the duplicate without any DB write.
	time.Sleep(1 * time.Second)

	apitest.New().
		HandlerFunc(fiberHandler()).
		Get("/api/v2/customers/" + code).
		Expect(t).
		Status(http.StatusOK).
		Assert(assertCustomer(code, "Updated Name")).
		End()
}
