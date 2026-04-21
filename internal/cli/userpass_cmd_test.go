package cli

import (
	"testing"
)

func TestParseUserpassFlags_Defaults(t *testing.T) {
	f, err := ParseUserpassFlags([]string{"-username", "alice", "-password", "secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Mount != "userpass" {
		t.Errorf("expected default mount 'userpass', got %q", f.Mount)
	}
	if f.Verbose {
		t.Error("expected verbose to be false by default")
	}
}

func TestParseUserpassFlags_AllFlags(t *testing.T) {
	f, err := ParseUserpassFlags([]string{
		"-address", "http://vault:8200",
		"-mount", "ldap",
		"-username", "bob",
		"-password", "hunter2",
		"-verbose",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Address != "http://vault:8200" {
		t.Errorf("unexpected address: %q", f.Address)
	}
	if f.Mount != "ldap" {
		t.Errorf("unexpected mount: %q", f.Mount)
	}
	if f.Username != "bob" {
		t.Errorf("unexpected username: %q", f.Username)
	}
	if f.Password != "hunter2" {
		t.Errorf("unexpected password: %q", f.Password)
	}
	if !f.Verbose {
		t.Error("expected verbose to be true")
	}
}

func TestParseUserpassFlags_MissingUsername(t *testing.T) {
	_, err := ParseUserpassFlags([]string{"-password", "secret"})
	if err == nil {
		t.Fatal("expected error for missing username")
	}
}

func TestParseUserpassFlags_MissingPassword(t *testing.T) {
	_, err := ParseUserpassFlags([]string{"-username", "alice"})
	if err == nil {
		t.Fatal("expected error for missing password")
	}
}

func TestParseUserpassFlags_InvalidFlag(t *testing.T) {
	_, err := ParseUserpassFlags([]string{"-unknown", "val"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}
