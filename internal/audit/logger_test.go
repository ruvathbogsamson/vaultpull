package audit_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/audit"
)

func TestLogger_Log_BasicEntry(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	e := audit.Entry{
		Operation:   "sync",
		SecretPath:  "secret/data/app",
		KeysWritten: 3,
		DryRun:      false,
	}

	if err := l.Log(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if got.Operation != "sync" {
		t.Errorf("expected operation sync, got %s", got.Operation)
	}
	if got.KeysWritten != 3 {
		t.Errorf("expected 3 keys written, got %d", got.KeysWritten)
	}
	if got.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLogger_Log_WithError(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	e := audit.Entry{
		Operation:  "sync",
		SecretPath: "secret/data/missing",
		Error:      "secret not found",
	}

	if err := l.Log(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if got.Error != "secret not found" {
		t.Errorf("expected error message, got %q", got.Error)
	}
}

func TestLogger_Log_TimestampPreserved(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	fixed := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	e := audit.Entry{
		Timestamp: fixed,
		Operation: "sync",
	}

	_ = l.Log(e)

	var got audit.Entry
	_ = json.Unmarshal(buf.Bytes(), &got)

	if !got.Timestamp.Equal(fixed) {
		t.Errorf("expected timestamp %v, got %v", fixed, got.Timestamp)
	}
}
