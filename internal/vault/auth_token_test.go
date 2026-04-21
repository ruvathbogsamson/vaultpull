package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newAuthTokenTestServer(status int, payload interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestNewAuthTokenClient_MissingAddress(t *testing.T) {
	_, err := NewAuthTokenClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewAuthTokenClient_MissingToken(t *testing.T) {
	_, err := NewAuthTokenClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestValidateToken_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"accessor":    "abc123",
			"policies":    []string{"default", "admin"},
			"ttl":         3600,
			"renewable":   true,
			"expire_time": "2099-01-01T00:00:00Z",
		},
	}
	srv := newAuthTokenTestServer(http.StatusOK, payload)
	defer srv.Close()

	c, err := NewAuthTokenClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := c.ValidateToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Accessor != "abc123" {
		t.Errorf("expected accessor abc123, got %s", info.Accessor)
	}
	if len(info.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.Policies))
	}
	if !info.Renewable {
		t.Error("expected token to be renewable")
	}
}

func TestValidateToken_Forbidden(t *testing.T) {
	srv := newAuthTokenTestServer(http.StatusForbidden, nil)
	defer srv.Close()

	c, _ := NewAuthTokenClient(srv.URL, "bad-token")
	_, err := c.ValidateToken()
	if err == nil {
		t.Fatal("expected error for forbidden token")
	}
}

func TestValidateToken_Unreachable(t *testing.T) {
	c, _ := NewAuthTokenClient("http://127.0.0.1:19999", "token")
	_, err := c.ValidateToken()
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
