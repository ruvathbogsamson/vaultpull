package audit

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newTestShipper(t *testing.T) (*Shipper, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	cfg := DefaultRotationConfig()
	s, err := NewShipper(path, cfg)
	if err != nil {
		t.Fatalf("NewShipper: %v", err)
	}
	return s, path
}

func TestShipper_Ship_WritesEntry(t *testing.T) {
	s, path := newTestShipper(t)

	e := Entry{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Operation: "sync",
		Path:      "secret/data/app",
		Namespace: "APP",
		Keys:      []string{"APP_DB_URL", "APP_SECRET"},
	}

	if err := s.Ship(e); err != nil {
		t.Fatalf("Ship: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open log: %v", err)
	}
	defer f.Close()

	var got Entry
	if err := json.NewDecoder(f).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got.Operation != "sync" {
		t.Errorf("operation = %q, want %q", got.Operation, "sync")
	}
	if got.Namespace != "APP" {
		t.Errorf("namespace = %q, want %q", got.Namespace, "APP")
	}
	if len(got.Keys) != 2 {
		t.Errorf("keys len = %d, want 2", len(got.Keys))
	}
}

func TestShipper_Ship_MultipleEntries(t *testing.T) {
	s, path := newTestShipper(t)

	for i := 0; i < 3; i++ {
		if err := s.Ship(Entry{Operation: "sync", Path: "secret/data/app"}); err != nil {
			t.Fatalf("Ship[%d]: %v", i, err)
		}
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() != "" {
			count++
		}
	}
	if count != 3 {
		t.Errorf("line count = %d, want 3", count)
	}
}

func TestShipper_Ship_SetsTimestamp(t *testing.T) {
	s, path := newTestShipper(t)

	if err := s.Ship(Entry{Operation: "fetch"}); err != nil {
		t.Fatalf("Ship: %v", err)
	}

	f, _ := os.Open(path)
	defer f.Close()
	var got Entry
	json.NewDecoder(f).Decode(&got) //nolint

	if got.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
