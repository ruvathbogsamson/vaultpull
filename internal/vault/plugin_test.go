package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newPluginTestServer(t *testing.T, pluginType string, plugins []PluginInfo, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(status)
		if status == http.StatusOK {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"detailed": plugins,
				},
			})
		}
	}))
}

func TestNewPluginClient_MissingAddress(t *testing.T) {
	_, err := NewPluginClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewPluginClient_MissingToken(t *testing.T) {
	_, err := NewPluginClient("http://localhost:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestListPlugins_Success(t *testing.T) {
	expected := []PluginInfo{
		{Name: "aws", Type: "secret", Version: "v1.0.0", Builtin: true},
		{Name: "kv", Type: "secret", Version: "v2.0.0", Builtin: true},
	}
	srv := newPluginTestServer(t, "secret", expected, http.StatusOK)
	defer srv.Close()

	client, err := NewPluginClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	plugins, err := client.ListPlugins("secret")
	if err != nil {
		t.Fatalf("ListPlugins failed: %v", err)
	}
	if len(plugins) != 2 {
		t.Fatalf("expected 2 plugins, got %d", len(plugins))
	}
	if plugins[0].Name != "aws" {
		t.Errorf("expected first plugin name 'aws', got %q", plugins[0].Name)
	}
}

func TestListPlugins_NotFound(t *testing.T) {
	srv := newPluginTestServer(t, "", nil, http.StatusNotFound)
	defer srv.Close()

	client, _ := NewPluginClient(srv.URL, "test-token")
	_, err := client.ListPlugins("unknown")
	if err == nil {
		t.Fatal("expected error for not found")
	}
}

func TestListPlugins_ServerError(t *testing.T) {
	srv := newPluginTestServer(t, "", nil, http.StatusInternalServerError)
	defer srv.Close()

	client, _ := NewPluginClient(srv.URL, "test-token")
	_, err := client.ListPlugins("secret")
	if err == nil {
		t.Fatal("expected error for server error")
	}
}
