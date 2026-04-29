package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/vaultpull/internal/vault"
)

func newControlGroupTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/sys/control-group/authorize":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"approved": true},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/v1/sys/control-group/request":
			var body map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["accessor"] == "" || body["accessor"] == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"approved":  false,
					"accessor":  body["accessor"],
					"request_id": "req-abc-123",
				},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestNewControlGroupClient_MissingAddress(t *testing.T) {
	_, err := vault.NewControlGroupClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewControlGroupClient_MissingToken(t *testing.T) {
	_, err := vault.NewControlGroupClient("http://localhost:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestControlGroupClient_Authorize_Success(t *testing.T) {
	srv := newControlGroupTestServer(t)
	defer srv.Close()

	client, err := vault.NewControlGroupClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ok, err := client.Authorize("acc-xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected authorized=true")
	}
}

func TestControlGroupClient_CheckRequest_Success(t *testing.T) {
	srv := newControlGroupTestServer(t)
	defer srv.Close()

	client, err := vault.NewControlGroupClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	status, err := client.CheckRequest("acc-xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.RequestID != "req-abc-123" {
		t.Errorf("expected request_id req-abc-123, got %s", status.RequestID)
	}
	if status.Approved {
		t.Error("expected approved=false")
	}
}

func TestControlGroupClient_CheckRequest_MissingAccessor(t *testing.T) {
	srv := newControlGroupTestServer(t)
	defer srv.Close()

	client, err := vault.NewControlGroupClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.CheckRequest("")
	if err == nil {
		t.Fatal("expected error for empty accessor")
	}
}
