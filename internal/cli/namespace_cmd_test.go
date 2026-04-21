package cli

import (
	"os"
	"testing"
)

func TestParseNamespaceFlags_Defaults(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")
	f, err := ParseNamespaceFlags([]string{"-namespace", "team"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Address != "http://127.0.0.1:8200" {
		t.Errorf("expected default address, got %q", f.Address)
	}
	if f.Namespace != "team" {
		t.Errorf("expected namespace 'team', got %q", f.Namespace)
	}
}

func TestParseNamespaceFlags_AllFlags(t *testing.T) {
	args := []string{
		"-addr", "https://vault.example.com",
		"-token", "s.abc123",
		"-namespace", "team/project",
	}
	f, err := ParseNamespaceFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Address != "https://vault.example.com" {
		t.Errorf("unexpected address: %q", f.Address)
	}
	if f.Token != "s.abc123" {
		t.Errorf("unexpected token: %q", f.Token)
	}
	if f.Namespace != "team/project" {
		t.Errorf("unexpected namespace: %q", f.Namespace)
	}
}

func TestParseNamespaceFlags_TokenFromEnv(t *testing.T) {
	os.Setenv("VAULT_TOKEN", "env-token")
	defer os.Unsetenv("VAULT_TOKEN")

	f, err := ParseNamespaceFlags([]string{"-namespace", "ops"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Token != "env-token" {
		t.Errorf("expected token from env, got %q", f.Token)
	}
}

func TestParseNamespaceFlags_MissingNamespace(t *testing.T) {
	_, err := ParseNamespaceFlags([]string{})
	if err == nil {
		t.Fatal("expected error for missing -namespace flag")
	}
}

func TestParseNamespaceFlags_InvalidFlag(t *testing.T) {
	_, err := ParseNamespaceFlags([]string{"-unknown", "val"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}
