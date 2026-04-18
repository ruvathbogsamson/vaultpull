package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempLog(t *testing.T, size int) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "audit-*.log")
	if err != nil {
		t.Fatalf("create temp log: %v", err)
	}
	defer f.Close()
	if size > 0 {
		data := make([]byte, size)
		if _, err := f.Write(data); err != nil {
			t.Fatalf("write temp log: %v", err)
		}
	}
	return f.Name()
}

func TestShouldRotate_SizeExceeded(t *testing.T) {
	path := writeTempLog(t, 1024)
	cfg := RotationConfig{MaxSizeBytes: 512, MaxAgeDays: 0}
	r := NewRotator(path, cfg)
	ok, err := r.ShouldRotate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected rotation due to size")
	}
}

func TestShouldRotate_BelowSize(t *testing.T) {
	path := writeTempLog(t, 100)
	cfg := RotationConfig{MaxSizeBytes: 10 * 1024 * 1024, MaxAgeDays: 0}
	r := NewRotator(path, cfg)
	ok, err := r.ShouldRotate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("did not expect rotation")
	}
}

func TestShouldRotate_NotExist(t *testing.T) {
	r := NewRotator("/nonexistent/audit.log", DefaultRotationConfig())
	ok, err := r.ShouldRotate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected no rotation for missing file")
	}
}

func TestRotate_CreatesBackup(t *testing.T) {
	dir := t.TempDir()
	path := writeTempLog(t, 256)
	cfg := RotationConfig{MaxSizeBytes: 1, MaxAgeDays: 0, BackupDir: dir}
	r := NewRotator(path, cfg)
	if err := r.Rotate(); err != nil {
		t.Fatalf("rotate: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("original log should be gone after rotation")
	}
	entries, err := filepath.Glob(filepath.Join(dir, "*.log.*"))
	if err != nil || len(entries) == 0 {
		t.Errorf("expected backup file in %s", dir)
	}
}

func TestDefaultRotationConfig(t *testing.T) {
	cfg := DefaultRotationConfig()
	if cfg.MaxSizeBytes != 10*1024*1024 {
		t.Errorf("unexpected MaxSizeBytes: %d", cfg.MaxSizeBytes)
	}
	if cfg.MaxAgeDays != 30 {
		t.Errorf("unexpected MaxAgeDays: %d", cfg.MaxAgeDays)
	}
	_ = time.Now() // ensure time import used
}
