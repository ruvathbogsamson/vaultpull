package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-org/vaultpull/internal/vault"
)

func newTOTPTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/totp/keys/myapp":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/totp/code/myapp":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]string{"code": "123456"},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/v1/totp/code/myapp":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]bool{"valid": true},
			})
		case r.URL.Path == "/v1/totp/keys/missing":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}

func TestNewTOTPClient_MissingAddress(t *testing.T) {
	_, err := vault.NewTOTPClient("", "token", "totp")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewTOTPClient_MissingToken(t *testing.T) {
	_, err := vault.NewTOTPClient("http://localhost", "", "totp")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewTOTPClient_DefaultMount(t *testing.T) {
	c, err := vault.NewTOTPClient("http://localhost", "token", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestGenerateCode_Success(t *testing.T) {
	srv := newTOTPTestServer(t)
	defer srv.Close()

	c, err := vault.NewTOTPClient(srv.URL, "test-token", "totp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	code, err := c.GenerateCode("myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != "123456" {
		t.Errorf("expected code 123456, got %s", code)
	}
}

func TestValidateCode_Success(t *testing.T) {
	srv := newTOTPTestServer(t)
	defer srv.Close()

	c, err := vault.NewTOTPClient(srv.URL, "test-token", "totp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	valid, err := c.ValidateCode("myapp", "123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !valid {
		t.Error("expected code to be valid")
	}
}

func TestGenerateCode_NotFound(t *testing.T) {
	srv := newTOTPTestServer(t)
	defer srv.Close()

	c, err := vault.NewTOTPClient(srv.URL, "test-token", "totp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = c.GenerateCode("missing")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}
