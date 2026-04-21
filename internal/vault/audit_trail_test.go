package vault

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAuditTrail_RecordAndEvents(t *testing.T) {
	trail := NewAuditTrail()
	trail.clock = fixedClock(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	trail.Record("read", "secret/app", "prod", []string{"DB_URL", "API_KEY"}, nil)

	events := trail.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	ev := events[0]
	if ev.Operation != "read" {
		t.Errorf("expected operation 'read', got %q", ev.Operation)
	}
	if ev.Path != "secret/app" {
		t.Errorf("expected path 'secret/app', got %q", ev.Path)
	}
	if ev.Namespace != "prod" {
		t.Errorf("expected namespace 'prod', got %q", ev.Namespace)
	}
	if len(ev.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(ev.Keys))
	}
	if ev.Error != "" {
		t.Errorf("expected no error, got %q", ev.Error)
	}
}

func TestAuditTrail_RecordWithError(t *testing.T) {
	trail := NewAuditTrail()
	trail.Record("write", "secret/db", "", nil, errors.New("permission denied"))

	events := trail.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Error != "permission denied" {
		t.Errorf("expected error 'permission denied', got %q", events[0].Error)
	}
}

func TestAuditTrail_Summary(t *testing.T) {
	trail := NewAuditTrail()
	if !strings.Contains(trail.Summary(), "no audit events") {
		t.Error("expected empty summary message")
	}

	trail.Record("read", "secret/app", "dev", []string{"KEY"}, nil)
	summary := trail.Summary()
	if !strings.Contains(summary, "read") || !strings.Contains(summary, "secret/app") {
		t.Errorf("summary missing expected content: %s", summary)
	}
}

func TestAuditTrail_Clear(t *testing.T) {
	trail := NewAuditTrail()
	trail.Record("read", "secret/app", "", nil, nil)
	trail.Clear()
	if len(trail.Events()) != 0 {
		t.Error("expected empty events after Clear")
	}
}

func TestAuditTrail_EventsIsolated(t *testing.T) {
	trail := NewAuditTrail()
	trail.Record("read", "secret/app", "", nil, nil)

	events := trail.Events()
	events[0].Operation = "mutated"

	original := trail.Events()
	if original[0].Operation != "read" {
		t.Error("Events() should return a copy, not a reference")
	}
}
