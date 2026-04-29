package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newEgressTestServer(t *testing.T, namespace string, rules []EgressRule, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(status)
		if status == http.StatusOK {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"rules": rules},
			})
		}
	}))
}

func TestNewEgressClient_MissingAddress(t *testing.T) {
	_, err := NewEgressClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewEgressClient_MissingToken(t *testing.T) {
	_, err := NewEgressClient("http://localhost:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestListRules_Success(t *testing.T) {
	rules := []EgressRule{
		{Path: "secret/app", Capabilities: []string{"read", "list"}},
		{Path: "secret/db", Capabilities: []string{"read"}},
	}
	srv := newEgressTestServer(t, "prod", rules, http.StatusOK)
	defer srv.Close()

	client, err := NewEgressClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := client.ListRules("prod")
	if err != nil {
		t.Fatalf("ListRules failed: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(got))
	}
	if got[0].Path != "secret/app" {
		t.Errorf("expected path secret/app, got %s", got[0].Path)
	}
}

func TestListRules_NotFound(t *testing.T) {
	srv := newEgressTestServer(t, "missing", nil, http.StatusNotFound)
	defer srv.Close()

	client, err := NewEgressClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.ListRules("missing")
	if err == nil {
		t.Fatal("expected error for not found namespace")
	}
}

func TestListRules_ServerError(t *testing.T) {
	srv := newEgressTestServer(t, "prod", nil, http.StatusInternalServerError)
	defer srv.Close()

	client, err := NewEgressClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.ListRules("prod")
	if err == nil {
		t.Fatal("expected error for server error")
	}
}
