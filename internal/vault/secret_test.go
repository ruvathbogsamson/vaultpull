package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newSecretTestServer(t *testing.T, path string, payload map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			t.Errorf("encoding response: %v", err)
		}
	}))
}

func TestFetchSecrets_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{
				"db_password": "s3cr3t",
				"api_key":     "abc123",
			},
		},
	}
	srv := newSecretTestServer(t, "/v1/secret/data/myapp", payload)
	defer srv.Close()

	c, err := New(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	secret, err := c.FetchSecrets("secret/myapp")
	if err != nil {
		t.Fatalf("FetchSecrets: %v", err)
	}

	if secret.Path != "secret/myapp" {
		t.Errorf("expected path %q, got %q", "secret/myapp", secret.Path)
	}

	if v, ok := secret.Data["DB_PASSWORD"]; !ok || v != "s3cr3t" {
		t.Errorf("expected DB_PASSWORD=s3cr3t, got %q", v)
	}
	if v, ok := secret.Data["API_KEY"]; !ok || v != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", v)
	}
}

func TestFetchSecrets_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c, err := New(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = c.FetchSecrets("secret/missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
