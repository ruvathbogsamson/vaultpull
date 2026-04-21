package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTokenTestServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/token/lookup-self" {
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

func TestLookupSelf_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"accessor":     "abc123",
			"policies":     []string{"default", "read-secrets"},
			"ttl":          3600,
			"renewable":    true,
			"display_name": "token-test",
			"meta":         map[string]string{"env": "staging"},
		},
	}
	srv := newTokenTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	client := NewTokenClient(srv.URL, "test-token", nil)
	info, err := client.LookupSelf()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Accessor != "abc123" {
		t.Errorf("expected accessor abc123, got %s", info.Accessor)
	}
	if len(info.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.Policies))
	}
	if info.TTL != 3600*time.Second {
		t.Errorf("expected TTL 3600s, got %v", info.TTL)
	}
	if !info.Renewable {
		t.Error("expected renewable to be true")
	}
	if info.DisplayName != "token-test" {
		t.Errorf("expected display_name token-test, got %s", info.DisplayName)
	}
}

func TestLookupSelf_Forbidden(t *testing.T) {
	srv := newTokenTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	client := NewTokenClient(srv.URL, "bad-token", nil)
	_, err := client.LookupSelf()
	if err == nil {
		t.Fatal("expected error for forbidden, got nil")
	}
}

func TestLookupSelf_Unreachable(t *testing.T) {
	client := NewTokenClient("http://127.0.0.1:19999", "tok", &http.Client{Timeout: 100 * time.Millisecond})
	_, err := client.LookupSelf()
	if err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
