package cli

import (
	"testing"
)

func TestParseTransitFlags_Defaults(t *testing.T) {
	f, err := ParseTransitFlags([]string{"-key", "mykey", "-payload", "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Mount != "transit" {
		t.Errorf("expected default mount 'transit', got %q", f.Mount)
	}
	if f.Action != "encrypt" {
		t.Errorf("expected default action 'encrypt', got %q", f.Action)
	}
}

func TestParseTransitFlags_AllFlags(t *testing.T) {
	f, err := ParseTransitFlags([]string{
		"-address", "http://vault:8200",
		"-token", "s.abc",
		"-mount", "my-transit",
		"-key", "aes256",
		"-action", "decrypt",
		"-payload", "vault:v1:xyz",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Address != "http://vault:8200" {
		t.Errorf("unexpected address: %q", f.Address)
	}
	if f.Mount != "my-transit" {
		t.Errorf("unexpected mount: %q", f.Mount)
	}
	if f.Action != "decrypt" {
		t.Errorf("unexpected action: %q", f.Action)
	}
}

func TestParseTransitFlags_MissingKey(t *testing.T) {
	_, err := ParseTransitFlags([]string{"-payload", "hello"})
	if err == nil {
		t.Fatal("expected error for missing -key")
	}
}

func TestParseTransitFlags_MissingPayload(t *testing.T) {
	_, err := ParseTransitFlags([]string{"-key", "mykey"})
	if err == nil {
		t.Fatal("expected error for missing -payload")
	}
}

func TestParseTransitFlags_InvalidAction(t *testing.T) {
	_, err := ParseTransitFlags([]string{"-key", "k", "-payload", "p", "-action", "sign"})
	if err == nil {
		t.Fatal("expected error for invalid action")
	}
}

func TestParseTransitFlags_InvalidFlag(t *testing.T) {
	_, err := ParseTransitFlags([]string{"-unknown"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}
