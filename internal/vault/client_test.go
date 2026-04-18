package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T, data map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"data": data,
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func TestNew_InvalidAddress(t *testing.T) {
	_, err := New("://bad-address", "token", "")
	if err == nil {
		t.Fatal("expected error for invalid address, got nil")
	}
}

func TestNew_ValidClient(t *testing.T) {
	c, err := New("http://127.0.0.1:8200", "test-token", "ns1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestMountFromPath(t *testing.T) {
	cases := []struct{ input, want string }{
		{"secret/myapp/prod", "secret"},
		{"kv/service", "kv"},
		{"onlymount", "onlymount"},
	}
	for _, tc := range cases {
		got := mountFromPath(tc.input)
		if got != tc.want {
			t.Errorf("mountFromPath(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestDataPathFromPath(t *testing.T) {
	cases := []struct{ input, want string }{
		{"secret/myapp/prod", "myapp/prod"},
		{"kv/service", "service"},
		{"onlymount", ""},
	}
	for _, tc := range cases {
		got := dataPathFromPath(tc.input)
		if got != tc.want {
			t.Errorf("dataPathFromPath(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestReadSecrets_NoData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"errors":["secret not found"]}`))
	}))
	defer server.Close()

	c, err := New(server.URL, "token", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = c.ReadSecrets(context.Background(), "secret/missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
