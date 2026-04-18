package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempLog(t *testing.T, size int) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "audit*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if size > 0 {
		if _, err := f.Write(make([]byte, size)); err != nil {
			t.Fatal(err)
		}
	}
	return f.Name()
}

func TestShouldRotate_SizeExceeded(t *testing.T) {
	path := writeTempLog(t, 11*1024*1024)
	r := NewRotator(path, DefaultRotationConfig())
	if !r.ShouldRotate() {
		t.Error("expected ShouldRotate true when size exceeded")
	}
}

func TestShouldRotate_BelowSize(t *testing.T) {
	path := writeTempLog(t, 100)
	r := NewRotator(path, DefaultRotationConfig())
	if r.ShouldRotate() {
		t.Error("expected ShouldRotate false when below threshold")
	}
}

func TestShouldRotate_NotExist(t *testing.T) {
	r := NewRotator("/nonexistent/audit.log", DefaultRotationConfig())
	if r.ShouldRotate() {
		t.Error("expected ShouldRotate false for missing file")
	}
}

func TestRotate_CreatesBackup(t *testing.T) {
	path := writeTempLog(t, 1024)
	r := NewRotator(path, DefaultRotationConfig())
	if err := r.Rotate(); err != nil {
		t.Fatalf("Rotate() error: %v", err)
	}
	// original should be gone
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected original log to be renamed away")
	}
	// backup should exist
	matches, _ := filepath.Glob(path + ".*")
	if len(matches) == 0 {
		t.Error("expected at least one backup file")
	}
}

func TestRotate_PrunesOldBackups(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "audit.log")
	cfg := RotationConfig{MaxSizeBytes: 1, MaxBackups: 2}
	r := NewRotator(base, cfg)

	// create 4 rotations
	for i := 0; i < 4; i++ {
		if err := os.WriteFile(base, []byte("x"), 0600); err != nil {
			t.Fatal(err)
		}
		if err := r.Rotate(); err != nil {
			t.Fatalf("Rotate() iteration %d: %v", i, err)
		}
	}
	matches, _ := filepath.Glob(base + ".*")
	if len(matches) > cfg.MaxBackups {
		t.Errorf("expected at most %d backups, got %d", cfg.MaxBackups, len(matches))
	}
}
