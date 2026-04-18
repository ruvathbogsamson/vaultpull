package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Path      string    `json:"path"`
	Namespace string    `json:"namespace,omitempty"`
	Keys      []string  `json:"keys,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// Shipper writes audit entries to a destination file as newline-delimited JSON.
type Shipper struct {
	path    string
	rotator *Rotator
}

// NewShipper creates a Shipper that appends entries to the given log file.
func NewShipper(path string, cfg RotationConfig) (*Shipper, error) {
	rotator, err := NewRotator(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("audit shipper: %w", err)
	}
	return &Shipper{path: path, rotator: rotator}, nil
}

// Ship encodes entry as JSON and appends it to the audit log.
func (s *Shipper) Ship(e Entry) error {
	if err := s.rotator.Rotate(); err != nil {
		return fmt.Errorf("audit shipper rotate: %w", err)
	}

	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("audit shipper open: %w", err)
	}
	defer f.Close()

	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	enc := json.NewEncoder(f)
	if err := enc.Encode(e); err != nil {
		return fmt.Errorf("audit shipper encode: %w", err)
	}
	return nil
}
