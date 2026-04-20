package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

func newLeaseTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/sys/renew":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"lease_id":       "secret/data/app#abc",
				"lease_duration": 3600,
				"renewable":      true,
			})
		case "/v1/sys/revoke":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func newLeaseClient(t *testing.T, addr string) *vaultapi.Client {
	t.Helper()
	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}
	client.SetToken("test-token")
	return client
}

func TestLeaseManager_TrackAndGet(t *testing.T) {
	server := newLeaseTestServer(t)
	defer server.Close()

	lm := NewLeaseManager(newLeaseClient(t, server.URL))
	lm.Track("secret/data/app#abc", 30*time.Minute, true)

	info := lm.Get("secret/data/app#abc")
	if info == nil {
		t.Fatal("expected lease info, got nil")
	}
	if info.Duration != 30*time.Minute {
		t.Errorf("expected 30m duration, got %v", info.Duration)
	}
	if !info.Renewable {
		t.Error("expected lease to be renewable")
	}
}

func TestLeaseManager_Renew(t *testing.T) {
	server := newLeaseTestServer(t)
	defer server.Close()

	lm := NewLeaseManager(newLeaseClient(t, server.URL))
	lm.Track("secret/data/app#abc", 30*time.Minute, true)

	if err := lm.Renew(context.Background(), "secret/data/app#abc", time.Hour); err != nil {
		t.Fatalf("unexpected renew error: %v", err)
	}

	info := lm.Get("secret/data/app#abc")
	if info.Duration != 3600*time.Second {
		t.Errorf("expected updated duration 3600s, got %v", info.Duration)
	}
}

func TestLeaseManager_Renew_NotTracked(t *testing.T) {
	server := newLeaseTestServer(t)
	defer server.Close()

	lm := NewLeaseManager(newLeaseClient(t, server.URL))
	err := lm.Renew(context.Background(), "nonexistent", time.Hour)
	if err == nil {
		t.Fatal("expected error for untracked lease")
	}
}

func TestLeaseManager_Expiring(t *testing.T) {
	server := newLeaseTestServer(t)
	defer server.Close()

	lm := NewLeaseManager(newLeaseClient(t, server.URL))
	lm.Track("lease/short", 5*time.Second, true)
	lm.Track("lease/long", 2*time.Hour, true)

	// manually backdate the short lease
	lm.leases["lease/short"].IssuedAt = time.Now().Add(-10 * time.Second)

	expiring := lm.Expiring(30 * time.Second)
	if len(expiring) != 1 {
		t.Fatalf("expected 1 expiring lease, got %d", len(expiring))
	}
	if expiring[0].LeaseID != "lease/short" {
		t.Errorf("expected lease/short, got %s", expiring[0].LeaseID)
	}
}
