package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newAuthMethodTestServer(status int, payload any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestNewAuthMethodClient_MissingAddress(t *testing.T) {
	_, err := NewAuthMethodClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewAuthMethodClient_MissingToken(t *testing.T) {
	_, err := NewAuthMethodClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestListAuthMethods_Success(t *testing.T) {
	payload := map[string]AuthMethod{
		"token/": {Type: "token", Description: "token based credentials", Accessor: "abc123", Local: false},
		"approle/": {Type: "approle", Description: "approle auth", Accessor: "def456", Local: true},
	}
	srv := newAuthMethodTestServer(http.StatusOK, payload)
	defer srv.Close()

	c, err := NewAuthMethodClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	methods, err := c.ListAuthMethods()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(methods) != 2 {
		t.Fatalf("expected 2 methods, got %d", len(methods))
	}
	if methods["token/"].Type != "token" {
		t.Errorf("expected token type, got %s", methods["token/"].Type)
	}
}

func TestListAuthMethods_Forbidden(t *testing.T) {
	srv := newAuthMethodTestServer(http.StatusForbidden, nil)
	defer srv.Close()

	c, _ := NewAuthMethodClient(srv.URL, "bad-token")
	_, err := c.ListAuthMethods()
	if err == nil {
		t.Fatal("expected permission denied error")
	}
}

func TestListAuthMethods_ServerError(t *testing.T) {
	srv := newAuthMethodTestServer(http.StatusInternalServerError, nil)
	defer srv.Close()

	c, _ := NewAuthMethodClient(srv.URL, "token")
	_, err := c.ListAuthMethods()
	if err == nil {
		t.Fatal("expected error on server error")
	}
}
