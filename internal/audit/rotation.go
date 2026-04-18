package audit

import (
	"fmt"
	"os"
	"time"
)

// RotationConfig holds configuration for log rotation.
type RotationConfig struct {
	MaxSizeBytes int64
	MaxBackups   int
}

// DefaultRotationConfig returns sensible defaults.
func DefaultRotationConfig() RotationConfig {
	return RotationConfig{
		MaxSizeBytes: 10 * 1024 * 1024, // 10 MB
		MaxBackups:   3,
	}
}

// Rotator manages log file rotation.
type Rotator struct {
	path   string
	cfg    RotationConfig
}

// NewRotator creates a new Rotator for the given log path.
func NewRotator(path string, cfg RotationConfig) *Rotator {
	return &Rotator{path: path, cfg: cfg}
}

// ShouldRotate reports whether the log file exceeds the max size.
func (r *Rotator) ShouldRotate() bool {
	info, err := os.Stat(r.path)
	if err != nil {
		return false
	}
	return info.Size() >= r.cfg.MaxSizeBytes
}

// Rotate renames the current log file to a timestamped backup and removes
// old backups beyond MaxBackups.
func (r *Rotator) Rotate() error {
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	backup := fmt.Sprintf("%s.%s", r.path, timestamp)
	if err := os.Rename(r.path, backup); err != nil {
		return fmt.Errorf("audit rotate: rename: %w", err)
	}
	return r.pruneBackups()
}

func (r *Rotator) pruneBackups() error {
	pattern := r.path + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("audit rotate: glob: %w", err)
	}
	// sort ascending; oldest first
	sort.Strings(matches)
	for len(matches) > r.cfg.MaxBackups {
		if err := os.Remove(matches[0]); err != nil {
			return fmt.Errorf("audit rotate: remove old backup: %w", err)
		}
		matches = matches[1:]
	}
	return nil
}
