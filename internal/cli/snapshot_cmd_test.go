package cli

import (
	"testing"
)

func TestParseSnapshotFlags_Defaults(t *testing.T) {
	flags, err := ParseSnapshotFlags([]string{"-path", "secret/app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.Address != "http://127.0.0.1:8200" {
		t.Errorf("expected default address, got %q", flags.Address)
	}
	if flags.OutputFile != "snapshot.json" {
		t.Errorf("expected default output file, got %q", flags.OutputFile)
	}
	if flags.Path != "secret/app" {
		t.Errorf("expected path secret/app, got %q", flags.Path)
	}
}

func TestFlags(t *testing.T) {
	flags, err := ParseSnapshotFlags([]string{
		"-addr", "http://vault:8200",
		"-token", "mytoken",
		"-path", "kv/prod",
		"-out", "prod.json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.Address != "http://vault:8200" {
		t.Errorf("address: want http://vault:8200, got %q", flags.Address)
	}
	if flags.Token != "mytoken" {
		t.Errorf("token: want mytoken, got %q", flags.Token)
	}
	if flags.Path != "kv/prod" {
		t.Errorf("path: want kv/prod, got %q", flags.Path)
	}
	if flags.OutputFile != "prod.json" {
		t.Errorf("out: want prod.json, got %q", flags.OutputFile)
	}
}

func TestParseSnapshotFlags_MissingPath(t *testing.T) {
	_, err := ParseSnapshotFlags([]string{"-addr", "http://vault:8200"})
	if err == nil {
		t.Fatal("expected error for missing -path, got nil")
	}
}

func TestParseSnapshotFlags_InvalidFlag(t *testing.T) {
	_, err := ParseSnapshotFlags([]string{"-unknown", "value", "-path", "secret/app"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}
