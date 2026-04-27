package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newInitTestServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/sys/init" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestNewInitClient_MissingAddress(t *testing.T) {
	_, err := NewInitClient("")
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestNewInitClient_Valid(t *testing.T) {
	c, err := NewInitClient("http://127.0.0.1:8200")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestInitialize_Success(t *testing.T) {
	expected := InitResponse{
		Keys:      []string{"key1", "key2", "key3"},
		RootToken: "s.roottoken",
	}
	srv := newInitTestServer(t, http.StatusOK, expected)
	defer srv.Close()

	c, _ := NewInitClient(srv.URL)
	resp, err := c.Initialize(InitRequest{SecretShares: 3, SecretThreshold: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.RootToken != expected.RootToken {
		t.Errorf("expected root token %q, got %q", expected.RootToken, resp.RootToken)
	}
	if len(resp.Keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(resp.Keys))
	}
}

func TestInitialize_InvalidShares(t *testing.T) {
	c, _ := NewInitClient("http://127.0.0.1:8200")
	_, err := c.Initialize(InitRequest{SecretShares: 0, SecretThreshold: 1})
	if err == nil {
		t.Fatal("expected error for zero secret_shares")
	}
}

func TestInitialize_InvalidThreshold(t *testing.T) {
	c, _ := NewInitClient("http://127.0.0.1:8200")
	_, err := c.Initialize(InitRequest{SecretShares: 3, SecretThreshold: 5})
	if err == nil {
		t.Fatal("expected error when threshold exceeds shares")
	}
}

func TestInitialize_ServerError(t *testing.T) {
	srv := newInitTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c, _ := NewInitClient(srv.URL)
	_, err := c.Initialize(InitRequest{SecretShares: 5, SecretThreshold: 3})
	if err == nil {
		t.Fatal("expected error on server error response")
	}
}
