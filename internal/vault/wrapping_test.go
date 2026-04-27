package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newWrappingTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/sys/wrapping/unwrap", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Vault-Token")
		if token == "bad-wrapping-token" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"DB_PASSWORD": "supersecret",
				"API_KEY":     "abc123",
			},
		})
	})

	mux.HandleFunc("/v1/sys/wrapping/lookup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"token":         "wrap-token-xyz",
				"accessor":      "acc-123",
				"creation_time": "2024-01-01T00:00:00Z",
			},
		})
	})

	return httptest.NewServer(mux)
}

func TestNewWrappingClient_MissingAddress(t *testing.T) {
	_, err := NewWrappingClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewWrappingClient_MissingToken(t *testing.T) {
	_, err := NewWrappingClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestUnwrap_Success(t *testing.T) {
	srv := newWrappingTestServer(t)
	defer srv.Close()

	client, err := NewWrappingClient(srv.URL, "valid-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := client.Unwrap("valid-wrapping-token")
	if err != nil {
		t.Fatalf("unexpected error unwrapping: %v", err)
	}
	if data["DB_PASSWORD"] != "supersecret" {
		t.Errorf("expected DB_PASSWORD=supersecret, got %q", data["DB_PASSWORD"])
	}
	if data["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", data["API_KEY"])
	}
}

func TestUnwrap_EmptyToken(t *testing.T) {
	srv := newWrappingTestServer(t)
	defer srv.Close()

	client, err := NewWrappingClient(srv.URL, "valid-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.Unwrap("")
	if err == nil {
		t.Fatal("expected error for empty wrapping token")
	}
}

func TestLookupWrappingToken_Success(t *testing.T) {
	srv := newWrappingTestServer(t)
	defer srv.Close()

	client, err := NewWrappingClient(srv.URL, "valid-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ws, err := client.LookupWrappingToken("wrap-token-xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ws.Accessor != "acc-123" {
		t.Errorf("expected accessor acc-123, got %q", ws.Accessor)
	}
	if ws.Creation != "2024-01-01T00:00:00Z" {
		t.Errorf("unexpected creation time: %q", ws.Creation)
	}
}
