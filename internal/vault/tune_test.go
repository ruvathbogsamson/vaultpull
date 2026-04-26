package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTuneTestServer(t *testing.T, mount string, cfg *TuneConfig, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/sys/mounts/" + mount + "/tune"
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if cfg != nil {
			_ = json.NewEncoder(w).Encode(cfg)
		}
	}))
}

func TestNewTuneClient_MissingAddress(t *testing.T) {
	_, err := NewTuneClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewTuneClient_MissingToken(t *testing.T) {
	_, err := NewTuneClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestGetTune_Success(t *testing.T) {
	cfg := &TuneConfig{
		DefaultLeaseTTL: "768h",
		MaxLeaseTTL:     "8760h",
		Description:     "kv secrets",
		ForceNoCache:    false,
	}
	srv := newTuneTestServer(t, "secret", cfg, http.StatusOK)
	defer srv.Close()

	client, err := NewTuneClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := client.GetTune("secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DefaultLeaseTTL != cfg.DefaultLeaseTTL {
		t.Errorf("expected DefaultLeaseTTL %q, got %q", cfg.DefaultLeaseTTL, result.DefaultLeaseTTL)
	}
	if result.MaxLeaseTTL != cfg.MaxLeaseTTL {
		t.Errorf("expected MaxLeaseTTL %q, got %q", cfg.MaxLeaseTTL, result.MaxLeaseTTL)
	}
}

func TestGetTune_NotFound(t *testing.T) {
	srv := newTuneTestServer(t, "secret", nil, http.StatusNotFound)
	defer srv.Close()

	client, err := NewTuneClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.GetTune("missing")
	if err == nil {
		t.Fatal("expected error for not found mount")
	}
}

func TestGetTune_MissingMount(t *testing.T) {
	client, err := NewTuneClient("http://127.0.0.1:8200", "token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = client.GetTune("")
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}
