package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newAliasCmdTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{
				"id":             "alias-abc",
				"name":           "devuser",
				"mount_accessor": "auth_ldap_xyz",
				"mount_type":     "ldap",
				"canonical_id":   "entity-999",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func TestParseAliasFlags_Defaults(t *testing.T) {
	flags, err := ParseAliasFlags([]string{"-id", "alias-abc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.AliasID != "alias-abc" {
		t.Errorf("expected alias-abc, got %s", flags.AliasID)
	}
}

func TestParseAliasFlags_AllFlags(t *testing.T) {
	flags, err := ParseAliasFlags([]string{
		"-address", "http://vault:8200",
		"-token", "mytoken",
		"-id", "alias-xyz",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.Address != "http://vault:8200" {
		t.Errorf("expected http://vault:8200, got %s", flags.Address)
	}
	if flags.Token != "mytoken" {
		t.Errorf("expected mytoken, got %s", flags.Token)
	}
}

func TestParseAliasFlags_MissingID(t *testing.T) {
	_, err := ParseAliasFlags([]string{"-address", "http://vault:8200"})
	if err == nil {
		t.Fatal("expected error for missing -id flag")
	}
}

func TestParseAliasFlags_InvalidFlag(t *testing.T) {
	_, err := ParseAliasFlags([]string{"-unknown", "value"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestRunAlias_Output(t *testing.T) {
	srv := newAliasCmdTestServer()
	defer srv.Close()

	flags := &AliasCmdFlags{
		Address: srv.URL,
		Token:   "test-token",
		AliasID: "alias-abc",
	}
	var buf bytes.Buffer
	if err := RunAlias(flags, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "alias-abc") {
		t.Errorf("expected alias-abc in output, got: %s", out)
	}
	if !strings.Contains(out, "devuser") {
		t.Errorf("expected devuser in output, got: %s", out)
	}
}
