package cli

import (
	"testing"
)

func TestParseEngineFlags_Defaults(t *testing.T) {
	flags, err := ParseEngineFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.Address != "http://127.0.0.1:8200" {
		t.Errorf("expected default address, got %q", flags.Address)
	}
	if flags.Mount != "" {
		t.Errorf("expected empty mount, got %q", flags.Mount)
	}
}

func TestParseEngineFlags_AllFlags(t *testing.T) {
	flags, err := ParseEngineFlags([]string{
		"-addr", "http://vault.example.com:8200",
		"-token", "s.mytoken",
		"-mount", "secret",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.Address != "http://vault.example.com:8200" {
		t.Errorf("expected custom address, got %q", flags.Address)
	}
	if flags.Token != "s.mytoken" {
		t.Errorf("expected token, got %q", flags.Token)
	}
	if flags.Mount != "secret" {
		t.Errorf("expected mount=secret, got %q", flags.Mount)
	}
}

func TestParseEngineFlags_TokenFromEnv(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "env-token")
	flags, err := ParseEngineFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.Token != "env-token" {
		t.Errorf("expected token from env, got %q", flags.Token)
	}
}

func TestParseEngineFlags_InvalidFlag(t *testing.T) {
	_, err := ParseEngineFlags([]string{"-unknown", "value"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}
