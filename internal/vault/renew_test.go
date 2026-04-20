package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

func newRenewTestServer(t *testing.T, renewCount *int32) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth/token/renew-self":
			atomic.AddInt32(renewCount, 1)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"auth": map[string]interface{}{"lease_duration": 3600},
			})
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestNewTokenRenewer_DefaultInterval(t *testing.T) {
	cfg := vaultapi.DefaultConfig()
	client, _ := vaultapi.NewClient(cfg)
	r := NewTokenRenewer(client, 0)
	if r.interval != 5*time.Minute {
		t.Errorf("expected default interval 5m, got %v", r.interval)
	}
}

func TestTokenRenewer_RenewsCalled(t *testing.T) {
	var count int32
	srv := newRenewTestServer(t, &count)
	defer srv.Close()

	cfg := vaultapi.DefaultConfig()
	cfg.Address = srv.URL
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}
	client.SetToken("test-token")

	renewer := NewTokenRenewer(client, 50*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	renewer.Start(ctx)
	<-ctx.Done()

	got := atomic.LoadInt32(&count)
	if got < 1 {
		t.Errorf("expected at least 1 renewal call, got %d", got)
	}
}

func TestTokenRenewer_Stop(t *testing.T) {
	var count int32
	srv := newRenewTestServer(t, &count)
	defer srv.Close()

	cfg := vaultapi.DefaultConfig()
	cfg.Address = srv.URL
	client, _ := vaultapi.NewClient(cfg)
	client.SetToken("test-token")

	renewer := NewTokenRenewer(client, 20*time.Millisecond)
	ctx := context.Background()
	renewer.Start(ctx)
	time.Sleep(30 * time.Millisecond)
	renewer.Stop()
	snap := atomic.LoadInt32(&count)
	time.Sleep(60 * time.Millisecond)
	after := atomic.LoadInt32(&count)
	if after > snap+1 {
		t.Errorf("renewals continued after Stop: before=%d after=%d", snap, after)
	}
}
