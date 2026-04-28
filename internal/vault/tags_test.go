package vault

import (
	"testing"
)

func TestNewTagFilter_Valid(t *testing.T) {
	f, err := NewTagFilter([]string{"env=prod", "team=backend"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Required["env"] != "prod" {
		t.Errorf("expected env=prod, got %s", f.Required["env"])
	}
	if f.Required["team"] != "backend" {
		t.Errorf("expected team=backend, got %s", f.Required["team"])
	}
}

func TestNewTagFilter_Invalid(t *testing.T) {
	cases := []string{"noequals", "=nokey", ""}
	for _, c := range cases {
		_, err := NewTagFilter([]string{c})
		if err == nil {
			t.Errorf("expected error for input %q, got nil", c)
		}
	}
}

func TestTagFilter_Match(t *testing.T) {
	f, _ := NewTagFilter([]string{"env=prod"})

	if !f.Match(map[string]string{"env": "prod", "region": "us-east"}) {
		t.Error("expected match")
	}
	if f.Match(map[string]string{"env": "staging"}) {
		t.Error("expected no match")
	}
	if f.Match(map[string]string{}) {
		t.Error("expected no match for missing key")
	}
}

func TestTagFilter_Match_Empty(t *testing.T) {
	f := &TagFilter{Required: map[string]string{}}
	if !f.Match(map[string]string{}) {
		t.Error("empty filter should match everything")
	}
}

func TestTagFilter_Match_MultipleRequired(t *testing.T) {
	f, _ := NewTagFilter([]string{"env=prod", "region=us-east"})

	// All required tags present and matching
	if !f.Match(map[string]string{"env": "prod", "region": "us-east"}) {
		t.Error("expected match when all required tags are present")
	}
	// Only one of the required tags matches
	if f.Match(map[string]string{"env": "prod", "region": "eu-west"}) {
		t.Error("expected no match when only one required tag matches")
	}
	// Neither required tag is present
	if f.Match(map[string]string{"team": "backend"}) {
		t.Error("expected no match when no required tags are present")
	}
}

func TestFilterSecrets_WithTags(t *testing.T) {
	secrets := map[string]string{
		"DB_PASS": "secret1",
		"API_KEY": "secret2",
		"TOKEN":   "secret3",
	}
	meta := map[string]map[string]string{
		"DB_PASS": {"env": "prod"},
		"API_KEY": {"env": "prod"},
		"TOKEN":   {"env": "staging"},
	}
	f, _ := NewTagFilter([]string{"env=prod"})
	out := FilterSecrets(secrets, meta, f)
	if len(out) != 2 {
		t.Fatalf("expected 2 secrets, got %d", len(out))
	}
	if _, ok := out["TOKEN"]; ok {
		t.Error("TOKEN should have been filtered out")
	}
}

func TestFilterSecrets_NilFilter(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	out := FilterSecrets(secrets, nil, nil)
	if len(out) != 2 {
		t.Errorf("nil filter should return all secrets, got %d", len(out))
	}
}
