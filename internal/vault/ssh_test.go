package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newSSHTestServer(t *testing.T, role, signedKey string, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/ssh/sign/" + role
		if r.URL.Path != expected {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(status)
		if status == http.StatusOK {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]string{
					"signed_key": signedKey,
					"key_type":   "ca",
				},
			})
		}
	}))
}

func TestNewSSHClient_MissingAddress(t *testing.T) {
	_, err := NewSSHClient("", "token", "")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewSSHClient_MissingToken(t *testing.T) {
	_, err := NewSSHClient("http://localhost", "", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewSSHClient_DefaultMount(t *testing.T) {
	c, err := NewSSHClient("http://localhost", "token", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.mount != "ssh" {
		t.Errorf("expected default mount 'ssh', got %q", c.mount)
	}
}

func TestSignKey_Success(t *testing.T) {
	const role = "developer"
	const signed = "ssh-rsa-cert-v01@openssh.com AAAA..."
	srv := newSSHTestServer(t, role, signed, http.StatusOK)
	defer srv.Close()

	c, _ := NewSSHClient(srv.URL, "test-token", "")
	cred, err := c.SignKey(role, "ssh-rsa AAAA...")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred.SignedKey != signed {
		t.Errorf("expected signed key %q, got %q", signed, cred.SignedKey)
	}
}

func TestSignKey_NotFound(t *testing.T) {
	srv := newSSHTestServer(t, "other", "", http.StatusNotFound)
	defer srv.Close()

	c, _ := NewSSHClient(srv.URL, "test-token", "")
	_, err := c.SignKey("missing", "ssh-rsa AAAA...")
	if err == nil {
		t.Fatal("expected error for missing role")
	}
}

func TestSignKey_MissingRole(t *testing.T) {
	c, _ := NewSSHClient("http://localhost", "token", "")
	_, err := c.SignKey("", "ssh-rsa AAAA...")
	if err == nil {
		t.Fatal("expected error for empty role")
	}
}

func TestSignKey_MissingPublicKey(t *testing.T) {
	c, _ := NewSSHClient("http://localhost", "token", "")
	_, err := c.SignKey("role", "")
	if err == nil {
		t.Fatal("expected error for empty public key")
	}
}
