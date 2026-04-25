package cli_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/cli"
)

func TestParseTOTPFlags_Defaults(t *testing.T) {
	f, err := cli.ParseTOTPFlags([]string{"-key", "myapp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Mount != "totp" {
		t.Errorf("expected default mount totp, got %s", f.Mount)
	}
	if f.Action != "generate" {
		t.Errorf("expected default action generate, got %s", f.Action)
	}
}

func TestParseTOTPFlags_AllFlags(t *testing.T) {
	f, err := cli.ParseTOTPFlags([]string{
		"-address", "http://vault:8200",
		"-token", "s.abc123",
		"-mount", "custom-totp",
		"-key", "myservice",
		"-action", "validate",
		"-code", "654321",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Address != "http://vault:8200" {
		t.Errorf("unexpected address: %s", f.Address)
	}
	if f.Mount != "custom-totp" {
		t.Errorf("unexpected mount: %s", f.Mount)
	}
	if f.KeyName != "myservice" {
		t.Errorf("unexpected key: %s", f.KeyName)
	}
	if f.Code != "654321" {
		t.Errorf("unexpected code: %s", f.Code)
	}
}

func TestParseTOTPFlags_MissingKey(t *testing.T) {
	_, err := cli.ParseTOTPFlags([]string{})
	if err == nil {
		t.Fatal("expected error for missing -key flag")
	}
}

func TestParseTOTPFlags_InvalidAction(t *testing.T) {
	_, err := cli.ParseTOTPFlags([]string{"-key", "myapp", "-action", "delete"})
	if err == nil {
		t.Fatal("expected error for invalid action")
	}
}

func TestParseTOTPFlags_ValidateMissingCode(t *testing.T) {
	_, err := cli.ParseTOTPFlags([]string{"-key", "myapp", "-action", "validate"})
	if err == nil {
		t.Fatal("expected error when validate action missing -code")
	}
}

func TestParseTOTPFlags_InvalidFlag(t *testing.T) {
	_, err := cli.ParseTOTPFlags([]string{"-unknown", "value"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}
