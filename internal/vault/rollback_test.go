package vault

import (
	"testing"
)

func TestRollbackManager_SaveAndLatest(t *testing.T) {
	rm := NewRollbackManager(3)

	secrets := map[string]string{"KEY": "value1"}
	rm.Save("secret/app", secrets)

	snap, err := rm.Latest("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Secrets["KEY"] != "value1" {
		t.Errorf("expected value1, got %s", snap.Secrets["KEY"])
	}
}

func TestRollbackManager_Previous(t *testing.T) {
	rm := NewRollbackManager(5)

	rm.Save("secret/app", map[string]string{"KEY": "v1"})
	rm.Save("secret/app", map[string]string{"KEY": "v2"})

	prev, err := rm.Previous("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prev.Secrets["KEY"] != "v1" {
		t.Errorf("expected v1, got %s", prev.Secrets["KEY"])
	}
}

func TestRollbackManager_Previous_NotEnoughSnapshots(t *testing.T) {
	rm := NewRollbackManager(5)
	rm.Save("secret/app", map[string]string{"KEY": "v1"})

	_, err := rm.Previous("secret/app")
	if err == nil {
		t.Fatal("expected error for missing previous snapshot")
	}
}

func TestRollbackManager_Latest_NoSnapshot(t *testing.T) {
	rm := NewRollbackManager(5)

	_, err := rm.Latest("secret/missing")
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestRollbackManager_MaxDepth(t *testing.T) {
	rm := NewRollbackManager(2)

	rm.Save("secret/app", map[string]string{"KEY": "v1"})
	rm.Save("secret/app", map[string]string{"KEY": "v2"})
	rm.Save("secret/app", map[string]string{"KEY": "v3"})

	// Only 2 snapshots should be retained; Previous should be v2.
	prev, err := rm.Previous("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prev.Secrets["KEY"] != "v2" {
		t.Errorf("expected v2 after eviction, got %s", prev.Secrets["KEY"])
	}
}

func TestRollbackManager_Flush(t *testing.T) {
	rm := NewRollbackManager(5)
	rm.Save("secret/app", map[string]string{"KEY": "v1"})
	rm.Flush("secret/app")

	_, err := rm.Latest("secret/app")
	if err == nil {
		t.Fatal("expected error after flush")
	}
}

func TestRollbackManager_IsolatesSnapshots(t *testing.T) {
	rm := NewRollbackManager(5)
	rm.Save("secret/a", map[string]string{"X": "1"})
	rm.Save("secret/b", map[string]string{"Y": "2"})

	snap, err := rm.Latest("secret/a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := snap.Secrets["Y"]; ok {
		t.Error("secret/a snapshot should not contain key Y from secret/b")
	}
}
