package cli_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/vaultpull/internal/cli"
)

func newControlGroupCmdTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/sys/control-group/authorize":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"approved": true},
			})
		case "/v1/sys/control-group/request":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"approved":   false,
					"accessor":   "acc-xyz",
					"request_id": "req-abc-123",
				},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestParseControlGroupFlags_Defaults(t *testing.T) {
	f, err := cli.ParseControlGroupFlags([]string{"--accessor", "acc-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Action != "check" {
		t.Errorf("expected default action 'check', got %q", f.Action)
	}
}

func TestParseControlGroupFlags_AllFlags(t *testing.T) {
	f, err := cli.ParseControlGroupFlags([]string{
		"--address", "http://vault:8200",
		"--token", "tok",
		"--accessor", "acc-xyz",
		"--action", "authorize",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Accessor != "acc-xyz" {
		t.Errorf("expected accessor acc-xyz, got %s", f.Accessor)
	}
	if f.Action != "authorize" {
		t.Errorf("expected action authorize, got %s", f.Action)
	}
}

func TestParseControlGroupFlags_MissingAccessor(t *testing.T) {
	_, err := cli.ParseControlGroupFlags([]string{"--action", "check"})
	if err == nil {
		t.Fatal("expected error for missing accessor")
	}
}

func TestParseControlGroupFlags_InvalidAction(t *testing.T) {
	_, err := cli.ParseControlGroupFlags([]string{"--accessor", "acc-1", "--action", "delete"})
	if err == nil {
		t.Fatal("expected error for invalid action")
	}
}

func TestRunControlGroup_Check(t *testing.T) {
	srv := newControlGroupCmdTestServer(t)
	defer srv.Close()

	f := cli.ControlGroupFlags{
		Address:  srv.URL,
		Token:    "test-token",
		Accessor: "acc-xyz",
		Action:   "check",
	}
	if err := cli.RunControlGroup(f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
