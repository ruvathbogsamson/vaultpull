package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newNamespaceTestServer(t *testing.T, keys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != "LIST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		body := map[string]interface{}{
			"data": map[string]interface{}{
				"keys": keys,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(body)
	}))
}

func TestNewNamespaceClient_MissingAddress(t *testing.T) {
	_, err := NewNamespaceClient("", "token", "team")
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestNewNamespaceClient_MissingNamespace(t *testing.T) {
	_, err := NewNamespaceClient("http://127.0.0.1:8200", "token", "")
	if err == nil {
		t.Fatal("expected error for empty namespace")
	}
}

func TestNewNamespaceClient_Valid(t *testing.T) {
	svr := newNamespaceTestServer(t, []string{"alpha/", "beta/"})
	defer svr.Close()

	nc, err := NewNamespaceClient(svr.URL, "test-token", "team")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nc.Namespace() != "team" {
		t.Errorf("expected namespace 'team', got %q", nc.Namespace())
	}
}

func TestNewNamespaceClient_StripsLeadingSlash(t *testing.T) {
	svr := newNamespaceTestServer(t, nil)
	defer svr.Close()

	nc, err := NewNamespaceClient(svr.URL, "token", "/team/project")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nc.Namespace() != "team/project" {
		t.Errorf("expected 'team/project', got %q", nc.Namespace())
	}
}

func TestListNamespaces_Empty(t *testing.T) {
	svr := newNamespaceTestServer(t, []string{})
	defer svr.Close()

	nc, err := NewNamespaceClient(svr.URL, "token", "team")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ns, err := nc.ListNamespaces()
	if err != nil {
		t.Fatalf("unexpected error listing: %v", err)
	}
	if len(ns) != 0 {
		t.Errorf("expected empty list, got %v", ns)
	}
}
