package env

import (
	"testing"
)

func TestFilter_EmptyNamespace(t *testing.T) {
	secrets := map[string]string{"APP_KEY": "1", "DB_PASS": "2"}
	result := Filter(secrets, "")
	if len(result) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(result))
	}
}

func TestFilter_WithNamespace(t *testing.T) {
	secrets := map[string]string{
		"APP_KEY":  "1",
		"APP_SECRET": "2",
		"DB_PASS":  "3",
	}
	result := Filter(secrets, "APP")
	if len(result) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(result))
	}
	if _, ok := result["DB_PASS"]; ok {
		t.Error("DB_PASS should have been filtered out")
	}
}

func TestFilter_CaseInsensitiveNamespace(t *testing.T) {
	secrets := map[string]string{
		"APP_KEY": "1",
		"DB_PASS": "2",
	}
	result := Filter(secrets, "app")
	if len(result) != 1 {
		t.Errorf("expected 1 secret, got %d", len(result))
	}
}

func TestFilter_NoMatch(t *testing.T) {
	secrets := map[string]string{"FOO_BAR": "1"}
	result := Filter(secrets, "NOMATCH")
	if len(result) != 0 {
		t.Errorf("expected 0 secrets, got %d", len(result))
	}
}
