package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/vault/api"
)

func newEngineTestServer(mounts map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/sys/mounts" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mounts)
			return
		}
		http.NotFound(w, r)
	}))
}

func newEngineClient(t *testing.T, srv *httptest.Server) *EngineClient {
	t.Helper()
	cfg := api.DefaultConfig()
	cfg.Address = srv.URL
	c, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("api.NewClient: %v", err)
	}
	c.SetToken("test-token")
	return NewEngineClient(c)
}

func TestListMounts_Success(t *testing.T) {
	srv := newEngineTestServer(map[string]interface{}{
		"secret/": map[string]interface{}{"type": "kv", "description": "KV store", "options": map[string]string{"version": "2"}},
		"sys/":    map[string]interface{}{"type": "system", "description": "System", "options": map[string]string{}},
	})
	defer srv.Close()

	ec := newEngineClient(t, srv)
	mounts, err := ec.ListMounts()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mounts) != 2 {
		t.Errorf("expected 2 mounts, got %d", len(mounts))
	}
}

func TestGetMount_Found(t *testing.T) {
	srv := newEngineTestServer(map[string]interface{}{
		"secret/": map[string]interface{}{"type": "kv", "description": "KV store", "options": map[string]string{"version": "2"}},
	})
	defer srv.Close()

	ec := newEngineClient(t, srv)
	m, err := ec.GetMount("secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Type != EngineKV1 {
		t.Errorf("expected type kv, got %s", m.Type)
	}
	if m.Version != "2" {
		t.Errorf("expected version 2, got %s", m.Version)
	}
}

func TestGetMount_NotFound(t *testing.T) {
	srv := newEngineTestServer(map[string]interface{}{
		"secret/": map[string]interface{}{"type": "kv", "description": "KV", "options": map[string]string{}},
	})
	defer srv.Close()

	ec := newEngineClient(t, srv)
	_, err := ec.GetMount("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing mount")
	}
}
