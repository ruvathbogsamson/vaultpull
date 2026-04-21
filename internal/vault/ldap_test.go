package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newLDAPTestServer(t *testing.T, statusCode int, token string, policies []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{
				"client_token": token,
				"policies":     policies,
			},
		})
	}))
}

func TestNewLDAPClient_MissingAddress(t *testing.T) {
	_, err := NewLDAPClient("", "")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewLDAPClient_DefaultMount(t *testing.T) {
	c, err := NewLDAPClient("http://127.0.0.1:8200", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.mount != "ldap" {
		t.Errorf("expected default mount 'ldap', got %q", c.mount)
	}
}

func TestLDAPLogin_Success(t *testing.T) {
	srv := newLDAPTestServer(t, http.StatusOK, "s.ldaptoken", []string{"default", "dev"})
	defer srv.Close()

	client, err := NewLDAPClient(srv.URL, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := client.Login("alice", "secret")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if resp.Token != "s.ldaptoken" {
		t.Errorf("expected token 's.ldaptoken', got %q", resp.Token)
	}
	if len(resp.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(resp.Policies))
	}
}

func TestLDAPLogin_MissingUsername(t *testing.T) {
	client, _ := NewLDAPClient("http://127.0.0.1:8200", "")
	_, err := client.Login("", "secret")
	if err == nil {
		t.Fatal("expected error for missing username")
	}
}

func TestLDAPLogin_MissingPassword(t *testing.T) {
	client, _ := NewLDAPClient("http://127.0.0.1:8200", "")
	_, err := client.Login("alice", "")
	if err == nil {
		t.Fatal("expected error for missing password")
	}
}

func TestLDAPLogin_Forbidden(t *testing.T) {
	srv := newLDAPTestServer(t, http.StatusForbidden, "", nil)
	defer srv.Close()

	client, _ := NewLDAPClient(srv.URL, "")
	_, err := client.Login("alice", "wrongpassword")
	if err == nil {
		t.Fatal("expected error for forbidden response")
	}
}
