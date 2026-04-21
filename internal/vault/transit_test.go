package vault

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTransitTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v1/transit/encrypt/mykey":
			var req map[string]string
			json.NewDecoder(r.Body).Decode(&req)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]string{"ciphertext": "vault:v1:abc123"},
			})
		case r.URL.Path == "/v1/transit/decrypt/mykey":
			w.Header().Set("Content-Type", "application/json")
			encoded := base64.StdEncoding.EncodeToString([]byte("hello"))
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]string{"plaintext": encoded},
			})
		case r.URL.Path == "/v1/transit/decrypt/badkey":
			w.WriteHeader(http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestNewTransitClient_MissingAddress(t *testing.T) {
	_, err := NewTransitClient("", "token", "")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewTransitClient_MissingToken(t *testing.T) {
	_, err := NewTransitClient("http://localhost", "", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewTransitClient_DefaultMount(t *testing.T) {
	c, err := NewTransitClient("http://localhost", "token", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.mount != "transit" {
		t.Errorf("expected default mount 'transit', got %q", c.mount)
	}
}

func TestTransitClient_Encrypt_Success(t *testing.T) {
	srv := newTransitTestServer(t)
	defer srv.Close()
	c, _ := NewTransitClient(srv.URL, "token", "transit")
	cipher, err := c.Encrypt("mykey", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cipher != "vault:v1:abc123" {
		t.Errorf("unexpected ciphertext: %q", cipher)
	}
}

func TestTransitClient_Decrypt_Success(t *testing.T) {
	srv := newTransitTestServer(t)
	defer srv.Close()
	c, _ := NewTransitClient(srv.URL, "token", "transit")
	plain, err := c.Decrypt("mykey", "vault:v1:abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plain != "hello" {
		t.Errorf("unexpected plaintext: %q", plain)
	}
}

func TestTransitClient_Decrypt_Error(t *testing.T) {
	srv := newTransitTestServer(t)
	defer srv.Close()
	c, _ := NewTransitClient(srv.URL, "token", "transit")
	_, err := c.Decrypt("badkey", "vault:v1:xyz")
	if err == nil {
		t.Fatal("expected error for forbidden key")
	}
}
