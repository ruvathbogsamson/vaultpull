package sync_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/sync"
)

func newVaultTestServer(t *testing.T, secrets map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"data": map[string]interface{}{
				"data": secrets,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}))
}

func TestSyncer_Run(t *testing.T) {
	srv := newVaultTestServer(t, map[string]string{
		"APP_KEY": "value1",
		"APP_SECRET": "value2",
		"OTHER_KEY": "value3",
	})
	defer srv.Close()

	out := filepath.Join(t.TempDir(), ".env")

	cfg := &config.Config{
		VaultAddr:  srv.URL,
		VaultToken: "test-token",
		SecretPath: "secret/data/myapp",
		Namespace:  "APP",
		OutputFile: out,
	}

	s, err := sync.New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	count, err := s.Run()
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 filtered secrets, got %d", count)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	content := string(data)
	for _, key := range []string{"APP_KEY", "APP_SECRET"} {
		if !containsKey(content, key) {
			t.Errorf("expected output to contain %q", key)
		}
	}
}

func containsKey(content, key string) bool {
	return len(content) > 0 && (len(key) == 0 || func() bool {
		for i := 0; i <= len(content)-len(key); i++ {
			if content[i:i+len(key)] == key {
				return true
			}
		}
		return false
	}())
}
