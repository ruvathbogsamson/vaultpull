package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriter_Write(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, ".env")

	w := NewWriter(path)
	secrets := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "BAZ=qux\nFOO=bar\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestWriter_Write_FilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, ".env")

	w := NewWriter(path)
	if err := w.Write(map[string]string{"KEY": "val"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}

	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
