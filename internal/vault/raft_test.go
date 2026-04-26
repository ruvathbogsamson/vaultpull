package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newRaftTestServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestNewRaftClient_MissingAddress(t *testing.T) {
	_, err := NewRaftClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewRaftClient_MissingToken(t *testing.T) {
	_, err := NewRaftClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestGetRaftConfig_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"index": 42,
			"servers": []map[string]interface{}{
				{"node_id": "node1", "address": "127.0.0.1:8201", "leader": true, "protocol_version": "3", "voter": true},
				{"node_id": "node2", "address": "127.0.0.1:8202", "leader": false, "protocol_version": "3", "voter": true},
			},
		},
	}
	srv := newRaftTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c, err := NewRaftClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, err := c.GetRaftConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Servers) != 2 {
		t.Errorf("expected 2 servers, got %d", len(cfg.Servers))
	}
	if !cfg.Servers[0].Leader {
		t.Error("expected first server to be leader")
	}
}

func TestGetRaftConfig_NotFound(t *testing.T) {
	srv := newRaftTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c, _ := NewRaftClient(srv.URL, "test-token")
	_, err := c.GetRaftConfig()
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestGetRaftConfig_ServerError(t *testing.T) {
	srv := newRaftTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c, _ := NewRaftClient(srv.URL, "test-token")
	_, err := c.GetRaftConfig()
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}
