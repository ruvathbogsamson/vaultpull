package cli

import (
	"testing"
)

func TestParseMFAFlags_Defaults(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://localhost:8200")
	t.Setenv("VAULT_TOKEN", "root")

	f, err := ParseMFAFlags([]string{"--request-id", "req-abc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Address != "http://localhost:8200" {
		t.Errorf("expected address from env, got %q", f.Address)
	}
	if f.Token != "root" {
		t.Errorf("expected token from env, got %q", f.Token)
	}
	if f.RequestID != "req-abc" {
		t.Errorf("expected request-id req-abc, got %q", f.RequestID)
	}
}

func TestParseMFAFlags_AllFlags(t *testing.T) {
	args := []string{
		"--address", "http://vault:8200",
		"--token", "mytoken",
		"--request-id", "req-xyz",
		"--payload", "654321",
	}
	f, err := ParseMFAFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Address != "http://vault:8200" {
		t.Errorf("unexpected address: %q", f.Address)
	}
	if f.Payload != "654321" {
		t.Errorf("unexpected payload: %q", f.Payload)
	}
}

func TestParseMFAFlags_MissingRequestID(t *testing.T) {
	_, err := ParseMFAFlags([]string{"--address", "http://localhost:8200"})
	if err == nil {
		t.Fatal("expected error for missing --request-id")
	}
}

func TestParseMFAFlags_InvalidFlag(t *testing.T) {
	_, err := ParseMFAFlags([]string{"--unknown-flag", "value"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}
