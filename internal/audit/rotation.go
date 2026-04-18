package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RotationConfig holds settings for log rotation.
type RotationConfig struct {
	MaxSizeBytes int64
	MaxAgeDays   int
	BackupDir    string
}

// DefaultRotationConfig returns sensible defaults.
func DefaultRotationConfig() RotationConfig {
	return RotationConfig{
		MaxSizeBytes: 10 * 1024 * 1024, // 10 MB
		MaxAgeDays:   30,
		BackupDir:    "",
	}
}

// Rotator manages audit log rotation.
type Rotator struct {
	cfg     RotationConfig
	logPath string
}

// NewRotator creates a Rotator for the given log file path.
func NewRotator(logPath string, cfg RotationConfig) *Rotator {
	return &Rotator{cfg: cfg, logPath: logPath}
}

// ShouldRotate reports whether the log file needs rotation.
func (r *Rotator) ShouldRotate() (bool, error) {
	info, err := os.Stat(r.logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("stat log file: %w", err)
	}
	if info.Size() >= r.cfg.MaxSizeBytes {
		return true, nil
	}
	if r.cfg.MaxAgeDays > 0 {
		cutoff := time.Now().AddDate(0, 0, -r.cfg.MaxAgeDays)
		if info.ModTime().Before(cutoff) {
			return true, nil
		}
	}
	return false, nil
}

// Rotate renames the current log file to a timestamped backup.
func (r *Rotator) Rotate() error {
	backupDir := r.cfg.BackupDir
	if backupDir == "" {
		backupDir = filepath.Dir(r.logPath)
	}
	if err := os.MkdirAll(backupDir, 0o700); err != nil {
		return fmt.Errorf("create backup dir: %w", err)
	}
	ts := time.Now().UTC().Format("20060102T150405Z")
	base := filepath.Base(r.logPath)
	dest := filepath.Join(backupDir, fmt.Sprintf("%s.%s", base, ts))
	if err := os.Rename(r.logPath, dest); err != nil {
		return fmt.Errorf("rotate log: %w", err)
	}
	return nil
}
