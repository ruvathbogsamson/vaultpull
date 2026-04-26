package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newAgentTestServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/agent/self" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestNewAgentClient_MissingAddress(t *testing.T) {
	_, err := NewAgentClient("")
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestNewAgentClient_Valid(t *testing.T) {
	client, err := NewAgentClient("http://127.0.0.1:8007")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestFetchConfig_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"token": "s.abc123",
			"mount": "auth/token",
			"metadata": map[string]string{"role": "web"},
		},
	}
	srv := newAgentTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	client, err := NewAgentClient(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, err := client.FetchConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Token != "s.abc123" {
		t.Errorf("expected token s.abc123, got %q", cfg.Token)
	}
	if cfg.Mount != "auth/token" {
		t.Errorf("expected mount auth/token, got %q", cfg.Mount)
	}
}

func TestFetchConfig_NotFound(t *testing.T) {
	srv := newAgentTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	client, _ := NewAgentClient(srv.URL)
	_, err := client.FetchConfig()
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestFetchConfig_Unreachable(t *testing.T) {
	client, _ := NewAgentClient("http://127.0.0.1:19999")
	_, err := client.FetchConfig()
	if err == nil {
		t.Fatal("expected error for unreachable agent")
	}
}
