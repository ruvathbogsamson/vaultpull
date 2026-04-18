package cli

import (
	"testing"
)

func TestParseFlags_Defaults(t *testing.T) {
	opts, err := ParseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.OutputFile != ".env" {
		t.Errorf("expected .env, got %s", opts.OutputFile)
	}
	if opts.Namespace != "" {
		t.Errorf("expected empty namespace, got %s", opts.Namespace)
	}
	if opts.DryRun {
		t.Error("expected dry-run to be false")
	}
	if opts.Verbose {
		t.Error("expected verbose to be false")
	}
}

func TestParseFlags_AllFlags(t *testing.T) {
	opts, err := ParseFlags([]string{
		"-output", "prod.env",
		"-namespace", "APP",
		"-dry-run",
		"-verbose",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.OutputFile != "prod.env" {
		t.Errorf("expected prod.env, got %s", opts.OutputFile)
	}
	if opts.Namespace != "APP" {
		t.Errorf("expected APP, got %s", opts.Namespace)
	}
	if !opts.DryRun {
		t.Error("expected dry-run to be true")
	}
	if !opts.Verbose {
		t.Error("expected verbose to be true")
	}
}

func TestParseFlags_InvalidFlag(t *testing.T) {
	_, err := ParseFlags([]string{"-unknown-flag"})
	if err == nil {
		t.Error("expected error for unknown flag")
	}
}
