// Package audit provides structured audit logging for vaultpull sync operations.
package audit

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp  time.Time `json:"timestamp"`
	Operation  string    `json:"operation"`
	SecretPath string    `json:"secret_path"`
	Namespace  string    `json:"namespace,omitempty"`
	KeysWritten int      `json:"keys_written"`
	DryRun     bool      `json:"dry_run"`
	Error      string    `json:"error,omitempty"`
}

// Logger writes structured audit entries as JSON lines.
type Logger struct {
	w io.Writer
}

// NewLogger returns a Logger writing to w. Pass nil to use os.Stdout.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{w: w}
}

// Log encodes e as a JSON line to the underlying writer.
func (l *Logger) Log(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = l.w.Write(b)
	return err
}
