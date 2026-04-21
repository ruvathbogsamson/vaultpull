package vault

import (
	"fmt"
	"time"
)

// AuditEvent represents a single secret access or modification event.
type AuditEvent struct {
	Timestamp  time.Time         `json:"timestamp"`
	Operation  string            `json:"operation"`
	Path       string            `json:"path"`
	Namespace  string            `json:"namespace,omitempty"`
	Keys       []string          `json:"keys,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Error      string            `json:"error,omitempty"`
}

// AuditTrail records vault operations for compliance and debugging.
type AuditTrail struct {
	events []AuditEvent
	clock  func() time.Time
}

// NewAuditTrail creates a new AuditTrail with the system clock.
func NewAuditTrail() *AuditTrail {
	return &AuditTrail{
		clock: time.Now,
	}
}

// Record appends an event to the trail.
func (a *AuditTrail) Record(op, path, namespace string, keys []string, err error) {
	ev := AuditEvent{
		Timestamp: a.clock(),
		Operation: op,
		Path:      path,
		Namespace: namespace,
		Keys:      keys,
	}
	if err != nil {
		ev.Error = err.Error()
	}
	a.events = append(a.events, ev)
}

// Events returns a copy of all recorded events.
func (a *AuditTrail) Events() []AuditEvent {
	out := make([]AuditEvent, len(a.events))
	copy(out, a.events)
	return out
}

// Summary returns a human-readable summary of all events.
func (a *AuditTrail) Summary() string {
	if len(a.events) == 0 {
		return "no audit events recorded"
	}
	summary := fmt.Sprintf("%d event(s) recorded:\n", len(a.events))
	for _, ev := range a.events {
		line := fmt.Sprintf("  [%s] %s %s", ev.Timestamp.Format(time.RFC3339), ev.Operation, ev.Path)
		if ev.Error != "" {
			line += " ERROR: " + ev.Error
		}
		summary += line + "\n"
	}
	return summary
}

// Clear removes all recorded events.
func (a *AuditTrail) Clear() {
	a.events = nil
}
