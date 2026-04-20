package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newPolicyTestServer(t *testing.T, name string, statusCode int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/sys/policy/" + name
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(statusCode)
		if statusCode == http.StatusOK {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"name":   name,
					"policy": `path "secret/*" { capabilities = ["read"] }`,
				},
			})
		}
	}))
}

func TestFetchPolicy_Success(t *testing.T) {
	srv := newPolicyTestServer(t, "readonly", http.StatusOK)
	defer srv.Close()

	c, err := New(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	policy, err := c.FetchPolicy(context.Background(), "readonly")
	if err != nil {
		t.Fatalf("FetchPolicy: %v", err)
	}
	if policy.Name != "readonly" {
		t.Errorf("expected policy name %q, got %q", "readonly", policy.Name)
	}
}

func TestFetchPolicy_NotFound(t *testing.T) {
	srv := newPolicyTestServer(t, "readonly", http.StatusNotFound)
	defer srv.Close()

	c, err := New(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = c.FetchPolicy(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error for missing policy, got nil")
	}
}

func TestFetchPolicy_ServerError(t *testing.T) {
	srv := newPolicyTestServer(t, "readonly", http.StatusInternalServerError)
	defer srv.Close()

	c, err := New(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = c.FetchPolicy(context.Background(), "readonly")
	if err == nil {
		t.Fatal("expected error for server error, got nil")
	}
}
