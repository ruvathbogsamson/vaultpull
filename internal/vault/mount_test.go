package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMountTestServer(t *testing.T, mounts map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/mounts" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mounts)
	}))
}

func TestNewMountClient_MissingAddress(t *testing.T) {
	_, err := NewMountClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewMountClient_MissingToken(t *testing.T) {
	_, err := NewMountClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestListMountPaths_Success(t *testing.T) {
	payload := map[string]interface{}{
		"secret/": map[string]interface{}{
			"type":        "kv",
			"description": "key/value store",
			"options":     map[string]string{"version": "2"},
		},
		"pki/": map[string]interface{}{
			"type":        "pki",
			"description": "PKI engine",
			"options":     map[string]string{},
		},
	}
	srv := newMountTestServer(t, payload)
	defer srv.Close()

	c, err := NewMountClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewMountClient: %v", err)
	}
	mounts, err := c.ListMountPaths()
	if err != nil {
		t.Fatalf("ListMountPaths: %v", err)
	}
	if len(mounts) != 2 {
		t.Fatalf("expected 2 mounts, got %d", len(mounts))
	}
	paths := map[string]bool{}
	for _, m := range mounts {
		paths[m.Path] = true
	}
	if !paths["secret"] {
		t.Error("expected mount path 'secret'")
	}
	if !paths["pki"] {
		t.Error("expected mount path 'pki'")
	}
}

func TestListMountPaths_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "forbidden", http.StatusForbidden)
	}))
	defer srv.Close()

	c, err := NewMountClient(srv.URL, "bad-token")
	if err != nil {
		t.Fatalf("NewMountClient: %v", err)
	}
	_, err = c.ListMountPaths()
	if err == nil {
		t.Fatal("expected error for forbidden response")
	}
}
