package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestParseWatchFlags_Defaults(t *testing.T) {
	var stderr bytes.Buffer
	f, err := ParseWatchFlags([]string{}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Interval != 60*time.Second {
		t.Errorf("expected default interval 60s, got %v", f.Interval)
	}
	if f.OutputFile != ".env" {
		t.Errorf("expected default output .env, got %q", f.OutputFile)
	}
	if f.Namespace != "" {
		t.Errorf("expected empty namespace, got %q", f.Namespace)
	}
	if f.Verbose {
		t.Error("expected verbose=false by default")
	}
}

func TestParseWatchFlags_AllFlags(t *testing.T) {
	var stderr bytes.Buffer
	args := []string{
		"-interval", "30s",
		"-output", "prod.env",
		"-namespace", "APP",
		"-verbose",
	}
	f, err := ParseWatchFlags(args, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Interval != 30*time.Second {
		t.Errorf("expected 30s, got %v", f.Interval)
	}
	if f.OutputFile != "prod.env" {
		t.Errorf("expected prod.env, got %q", f.OutputFile)
	}
	if f.Namespace != "APP" {
		t.Errorf("expected APP, got %q", f.Namespace)
	}
	if !f.Verbose {
		t.Error("expected verbose=true")
	}
}

func TestParseWatchFlags_InvalidFlag(t *testing.T) {
	var stderr bytes.Buffer
	_, err := ParseWatchFlags([]string{"-unknown"}, &stderr)
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
	if !strings.Contains(stderr.String(), "unknown") && !strings.Contains(err.Error(), "unknown") {
		t.Errorf("expected error message to mention unknown flag, got: %v / %s", err, stderr.String())
	}
}
