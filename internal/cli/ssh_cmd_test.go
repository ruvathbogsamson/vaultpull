package cli

import (
	"testing"
)

func TestParseSSHFlags_Defaults(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://vault:8200")
	t.Setenv("VAULT_TOKEN", "root")

	f, err := ParseSSHFlags([]string{"-role", "dev", "-public-key", "ssh-rsa AAAA"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Address != "http://vault:8200" {
		t.Errorf("expected address from env, got %q", f.Address)
	}
	if f.Token != "root" {
		t.Errorf("expected token from env, got %q", f.Token)
	}
	if f.Mount != "ssh" {
		t.Errorf("expected default mount 'ssh', got %q", f.Mount)
	}
}

func TestParseSSHFlags_AllFlags(t *testing.T) {
	f, err := ParseSSHFlags([]string{
		"-address", "http://localhost:8200",
		"-token", "mytoken",
		"-mount", "custom-ssh",
		"-role", "ops",
		"-public-key", "ssh-ed25519 AAAA",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Mount != "custom-ssh" {
		t.Errorf("expected mount 'custom-ssh', got %q", f.Mount)
	}
	if f.Role != "ops" {
		t.Errorf("expected role 'ops', got %q", f.Role)
	}
	if f.PublicKey != "ssh-ed25519 AAAA" {
		t.Errorf("unexpected public key: %q", f.PublicKey)
	}
}

func TestParseSSHFlags_MissingRole(t *testing.T) {
	_, err := ParseSSHFlags([]string{"-public-key", "ssh-rsa AAAA"})
	if err == nil {
		t.Fatal("expected error for missing -role")
	}
}

func TestParseSSHFlags_MissingPublicKey(t *testing.T) {
	_, err := ParseSSHFlags([]string{"-role", "dev"})
	if err == nil {
		t.Fatal("expected error for missing -public-key")
	}
}

func TestParseSSHFlags_InvalidFlag(t *testing.T) {
	_, err := ParseSSHFlags([]string{"-unknown", "val"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}
