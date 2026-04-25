package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newSealTestServer(t *testing.T, sealed bool) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/sys/seal-status":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(SealStatus{
				Sealed:      sealed,
				Initialized: true,
				Version:     "1.15.0",
				ClusterName: "vault-cluster",
			})
		case "/v1/sys/seal":
			if r.Method != http.MethodPut {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestNewSealClient_MissingAddress(t *testing.T) {
	_, err := NewSealClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewSealClient_MissingToken(t *testing.T) {
	_, err := NewSealClient("http://localhost:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestGetSealStatus_Unsealed(t *testing.T) {
	srv := newSealTestServer(t, false)
	defer srv.Close()

	client, err := NewSealClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	status, err := client.GetSealStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Sealed {
		t.Error("expected vault to be unsealed")
	}
	if !status.Initialized {
		t.Error("expected vault to be initialized")
	}
	if status.Version != "1.15.0" {
		t.Errorf("expected version 1.15.0, got %s", status.Version)
	}
}

func TestGetSealStatus_Sealed(t *testing.T) {
	srv := newSealTestServer(t, true)
	defer srv.Close()

	client, err := NewSealClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	status, err := client.GetSealStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Sealed {
		t.Error("expected vault to be sealed")
	}
}

func TestSeal_Success(t *testing.T) {
	srv := newSealTestServer(t, false)
	defer srv.Close()

	client, err := NewSealClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.Seal(); err != nil {
		t.Fatalf("unexpected error sealing: %v", err)
	}
}

func TestGetSealStatus_Unreachable(t *testing.T) {
	client, _ := NewSealClient("http://127.0.0.1:19999", "token")
	_, err := client.GetSealStatus()
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
