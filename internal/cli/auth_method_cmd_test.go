package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/vaultpull/internal/vault"
)

func newAuthMethodCmdTestServer() *httptest.Server {
	payload := map[string]vault.AuthMethod{
		"token/":   {Type: "token", Description: "token auth", Accessor: "a1", Local: false},
		"approle/": {Type: "approle", Description: "approle auth", Accessor: "a2", Local: true},
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func TestParseAuthMethodFlags_Defaults(t *testing.T) {
	f, err := ParseAuthMethodFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil flags")
	}
}

func TestParseAuthMethodFlags_AllFlags(t *testing.T) {
	f, err := ParseAuthMethodFlags([]string{"-address", "http://vault:8200", "-token", "mytoken"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Address != "http://vault:8200" {
		t.Errorf("expected address http://vault:8200, got %s", f.Address)
	}
	if f.Token != "mytoken" {
		t.Errorf("expected token mytoken, got %s", f.Token)
	}
}

func TestParseAuthMethodFlags_InvalidFlag(t *testing.T) {
	_, err := ParseAuthMethodFlags([]string{"-unknown"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestRunAuthMethod_Output(t *testing.T) {
	srv := newAuthMethodCmdTestServer()
	defer srv.Close()

	f := &AuthMethodFlags{Address: srv.URL, Token: "test-token"}
	var buf bytes.Buffer
	if err := RunAuthMethod(f, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "token/") {
		t.Errorf("expected token/ in output, got: %s", out)
	}
	if !strings.Contains(out, "approle/") {
		t.Errorf("expected approle/ in output, got: %s", out)
	}
}
