package vault

import (
	"strings"
	"testing"
)

func TestDiff_AllAdded(t *testing.T) {
	remote := map[string]string{"FOO": "bar", "BAZ": "qux"}
	local := map[string]string{}

	d := Diff(remote, local)

	if len(d.Added) != 2 {
		t.Fatalf("expected 2 added, got %d", len(d.Added))
	}
	if d.HasChanges() == false {
		t.Fatal("expected HasChanges to be true")
	}
}

func TestDiff_AllRemoved(t *testing.T) {
	remote := map[string]string{}
	local := map[string]string{"FOO": "bar", "BAZ": "qux"}

	d := Diff(remote, local)

	if len(d.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(d.Removed))
	}
}

func TestDiff_Changed(t *testing.T) {
	remote := map[string]string{"FOO": "newval"}
	local := map[string]string{"FOO": "oldval"}

	d := Diff(remote, local)

	if len(d.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(d.Changed))
	}
	if d.Changed["FOO"] != "newval" {
		t.Errorf("expected changed value 'newval', got %q", d.Changed["FOO"])
	}
}

func TestDiff_Unchanged(t *testing.T) {
	remote := map[string]string{"FOO": "bar"}
	local := map[string]string{"FOO": "bar"}

	d := Diff(remote, local)

	if len(d.Unchanged) != 1 {
		t.Fatalf("expected 1 unchanged, got %d", len(d.Unchanged))
	}
	if d.HasChanges() {
		t.Fatal("expected HasChanges to be false")
	}
}

func TestDiff_Mixed(t *testing.T) {
	remote := map[string]string{"ADDED": "v1", "CHANGED": "new", "SAME": "x"}
	local := map[string]string{"REMOVED": "v2", "CHANGED": "old", "SAME": "x"}

	d := Diff(remote, local)

	if len(d.Added) != 1 || d.Added["ADDED"] != "v1" {
		t.Errorf("unexpected Added: %v", d.Added)
	}
	if len(d.Removed) != 1 || d.Removed["REMOVED"] != "v2" {
		t.Errorf("unexpected Removed: %v", d.Removed)
	}
	if len(d.Changed) != 1 || d.Changed["CHANGED"] != "new" {
		t.Errorf("unexpected Changed: %v", d.Changed)
	}
	if len(d.Unchanged) != 1 {
		t.Errorf("unexpected Unchanged: %v", d.Unchanged)
	}
}

func TestDiff_Summary(t *testing.T) {
	remote := map[string]string{"A": "1", "B": "new"}
	local := map[string]string{"B": "old", "C": "3"}

	d := Diff(remote, local)
	s := d.Summary()

	if !strings.Contains(s, "+1") {
		t.Errorf("summary missing added count: %s", s)
	}
	if !strings.Contains(s, "~1") {
		t.Errorf("summary missing changed count: %s", s)
	}
	if !strings.Contains(s, "-1") {
		t.Errorf("summary missing removed count: %s", s)
	}
}
