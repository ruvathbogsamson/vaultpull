package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newUserpassTestServer(t *testing.T, statusCode int, token string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if statusCode == http.StatusOK {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"auth": map[string]interface{}{
					"client_token":   token,
					"lease_duration": 3600,
					"renewable":      true,
				},
			})
		}
	}))
}

func TestNewUserpassClient_MissingAddress(t *testing.T) {
	_, err := NewUserpassClient("", "")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewUserpassClient_DefaultMount(t *testing.T) {
	c, err := NewUserpassClient("http://127.0.0.1:8200", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.mount != "userpass" {
		t.Errorf("expected default mount 'userpass', got %q", c.mount)
	}
}

func TestUserpassLogin_Success(t *testing.T) {
	srv := newUserpassTestServer(t, http.StatusOK, "s.testtoken")
	defer srv.Close()

	c, _ := NewUserpassClient(srv.URL, "")
	tok, err := c.Login(UserpassCredentials{Username: "alice", Password: "secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.ClientToken != "s.testtoken" {
		t.Errorf("expected token 's.testtoken', got %q", tok.ClientToken)
	}
	if tok.LeaseDuration != 3600 {
		t.Errorf("expected lease 3600, got %d", tok.LeaseDuration)
	}
	if !tok.Renewable {
		t.Error("expected renewable to be true")
	}
}

func TestUserpassLogin_MissingUsername(t *testing.T) {
	c, _ := NewUserpassClient("http://127.0.0.1:8200", "")
	_, err := c.Login(UserpassCredentials{Username: "", Password: "secret"})
	if err == nil {
		t.Fatal("expected error for missing username")
	}
}

func TestUserpassLogin_MissingPassword(t *testing.T) {
	c, _ := NewUserpassClient("http://127.0.0.1:8200", "")
	_, err := c.Login(UserpassCredentials{Username: "alice", Password: ""})
	if err == nil {
		t.Fatal("expected error for missing password")
	}
}

func TestUserpassLogin_Forbidden(t *testing.T) {
	srv := newUserpassTestServer(t, http.StatusForbidden, "")
	defer srv.Close()

	c, _ := NewUserpassClient(srv.URL, "")
	_, err := c.Login(UserpassCredentials{Username: "alice", Password: "wrong"})
	if err == nil {
		t.Fatal("expected error for forbidden")
	}
}
