package cli

import (
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/vault"
)

func newTestTrail() *vault.AuditTrail {
	t := vault.NewAuditTrail()
	return t
}

func TestParseAuditTrailFlags_Defaults(t *testing.T) {
	flags, err := ParseAuditTrailFlags([]string{"-path", "secret/app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.Path != "secret/app" {
		t.Errorf("expected path 'secret/app', got %q", flags.Path)
	}
	if flags.Namespace != "" {
		t.Errorf("expected empty namespace, got %q", flags.Namespace)
	}
	if flags.Verbose {
		t.Error("expected verbose=false by default")
	}
}

func TestParseAuditTrailFlags_AllFlags(t *testing.T) {
	flags, err := ParseAuditTrailFlags([]string{
		"-path", "secret/db",
		"-namespace", "prod",
		"-verbose",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.Namespace != "prod" {
		t.Errorf("expected namespace 'prod', got %q", flags.Namespace)
	}
	if !flags.Verbose {
		t.Error("expected verbose=true")
	}
}

func TestParseAuditTrailFlags_MissingPath(t *testing.T) {
	_, err := ParseAuditTrailFlags([]string{})
	if err == nil {
		t.Fatal("expected error for missing -path flag")
	}
}

func TestRunAuditTrail_Empty(t *testing.T) {
	trail := newTestTrail()
	flags := &AuditTrailFlags{Path: "secret/app"}
	if err := RunAuditTrail(flags, trail); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAuditTrail_WithEvents(t *testing.T) {
	trail := newTestTrail()
	_ = time.Now()
	trail.Record("read", "secret/app", "prod", []string{"DB_URL"}, nil)
	trail.Record("read", "secret/app", "dev", []string{"API_KEY"}, nil)

	flags := &AuditTrailFlags{Path: "secret/app", Namespace: "prod"}
	if err := RunAuditTrail(flags, trail); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
