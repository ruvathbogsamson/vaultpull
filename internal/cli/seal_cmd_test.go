package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/vault"
)

func newSealCmdTestServer(t *testing.T, sealed bool) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/sys/seal-status":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(vault.SealStatus{
				Sealed:      sealed,
				Initialized: true,
				Version:     "1.15.0",
				ClusterName: "test-cluster",
			})
		case "/v1/sys/seal":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestParseSealFlags_Defaults(t *testing.T) {
	f, err := ParseSealFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Action != "status" {
		t.Errorf("expected default action 'status', got %q", f.Action)
	}
}

func TestParseSealFlags_AllFlags(t *testing.T) {
	f, err := ParseSealFlags([]string{
		"-address", "http://vault:8200",
		"-token", "s.abc",
		"-action", "seal",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Address != "http://vault:8200" {
		t.Errorf("unexpected address: %s", f.Address)
	}
	if f.Action != "seal" {
		t.Errorf("unexpected action: %s", f.Action)
	}
}

func TestParseSealFlags_InvalidAction(t *testing.T) {
	_, err := ParseSealFlags([]string{"-action", "unseal"})
	if err == nil {
		t.Fatal("expected error for invalid action")
	}
}

func TestRunSeal_Status_Unsealed(t *testing.T) {
	srv := newSealCmdTestServer(t, false)
	defer srv.Close()

	f := &SealFlags{Address: srv.URL, Token: "test-token", Action: "status"}
	var buf bytes.Buffer
	if err := RunSeal(f, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "unsealed") {
		t.Errorf("expected 'unsealed' in output, got: %s", buf.String())
	}
}

func TestRunSeal_Seal_Success(t *testing.T) {
	srv := newSealCmdTestServer(t, false)
	defer srv.Close()

	f := &SealFlags{Address: srv.URL, Token: "test-token", Action: "seal"}
	var buf bytes.Buffer
	if err := RunSeal(f, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "sealed successfully") {
		t.Errorf("expected success message, got: %s", buf.String())
	}
}
