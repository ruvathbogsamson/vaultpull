package cli

import (
	"testing"
)

func TestParseKVFlags_Defaults(t *testing.T) {
	f, err := ParseKVFlags([]string{"-path", "myapp/db"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Mount != "secret" {
		t.Errorf("Mount: got %q, want %q", f.Mount, "secret")
	}
	if f.KVVersion != 0 {
		t.Errorf("KVVersion: got %d, want 0", f.KVVersion)
	}
	if f.Verbose {
		t.Error("Verbose should default to false")
	}
}

func TestParseKVFlags_AllFlags(t *testing.T) {
	f, err := ParseKVFlags([]string{
		"-mount", "ops",
		"-path", "infra/redis",
		"-kv-version", "2",
		"-verbose",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Mount != "ops" {
		t.Errorf("Mount: got %q, want %q", f.Mount, "ops")
	}
	if f.SecretPath != "infra/redis" {
		t.Errorf("SecretPath: got %q, want %q", f.SecretPath, "infra/redis")
	}
	if f.KVVersion != 2 {
		t.Errorf("KVVersion: got %d, want 2", f.KVVersion)
	}
	if !f.Verbose {
		t.Error("Verbose should be true")
	}
}

func TestParseKVFlags_MissingPath(t *testing.T) {
	_, err := ParseKVFlags([]string{"-mount", "secret"})
	if err == nil {
		t.Fatal("expected error for missing -path")
	}
}

func TestParseKVFlags_InvalidKVVersion(t *testing.T) {
	_, err := ParseKVFlags([]string{"-path", "x", "-kv-version", "3"})
	if err == nil {
		t.Fatal("expected error for invalid -kv-version")
	}
}

func TestParseKVFlags_InvalidFlag(t *testing.T) {
	_, err := ParseKVFlags([]string{"-unknown-flag"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestParseKVFlags_KVVersion1(t *testing.T) {
	f, err := ParseKVFlags([]string{"-path", "myapp/config", "-kv-version", "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.KVVersion != 1 {
		t.Errorf("KVVersion: got %d, want 1", f.KVVersion)
	}
}
