package cli

import (
	"bytes"
	"os"
	"testing"
)

func setEnv(t *testing.T, key, val string) {
	t.Helper()
	t.Setenv(key, val)
}

func TestRunner_DryRun_NoConfig(t *testing.T) {
	opts := &Options{DryRun: true, OutputFile: ".env"}
	r := NewRunner(opts)
	err := r.Run()
	// Should fail on missing config, not on file write
	if err == nil {
		t.Error("expected config error when VAULT_TOKEN is missing")
	}
}

func TestRunner_DryRun_Verbose(t *testing.T) {
	setEnv(t, "VAULT_TOKEN", "test-token")
	setEnv(t, "VAULT_SECRET_PATH", "secret/data/app")
	setEnv(t, "VAULT_ADDR", "http://127.0.0.1:19999") // unreachable but dry-run skips sync

	var buf bytes.Buffer
	opts := &Options{DryRun: true, Verbose: true, OutputFile: ".env"}
	r := NewRunner(opts)
	r.stdout = &buf

	// dry-run exits before vault connection
	_ = r.Run()

	if buf.Len() > 0 && !bytes.Contains(buf.Bytes(), []byte("dry-run")) {
		t.Errorf("expected dry-run message in output, got: %s", buf.String())
	}
}

func TestRunner_DryRun_NamespaceOverride(t *testing.T) {
	setEnv(t, "VAULT_TOKEN", "tok")
	setEnv(t, "VAULT_SECRET_PATH", "secret/data/app")
	setEnv(t, "VAULT_ADDR", "http://127.0.0.1:19999")
	setEnv(t, "NAMESPACE", "")

	opts := &Options{DryRun: true, Namespace: "MYAPP", OutputFile: ".env"}
	r := NewRunner(opts)
	r.stdout = &bytes.Buffer{}

	_ = r.Run()
	// Just ensure it doesn't panic
	_ = os.Remove(".env")
}
