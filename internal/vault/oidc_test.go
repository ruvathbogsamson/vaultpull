package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newOIDCTestServer(t *testing.T, status int, token string, policies []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/oidc/login" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if status == http.StatusOK {
			body := map[string]interface{}{
				"auth": map[string]interface{}{
					"client_token":   token,
					"policies":       policies,
					"lease_duration": 3600,
				},
			}
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestNewOIDCClient_MissingAddress(t *testing.T) {
	_, err := NewOIDCClient("", "myrole")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewOIDCClient_MissingRole(t *testing.T) {
	_, err := NewOIDCClient("http://localhost:8200", "")
	if err == nil {
		t.Fatal("expected error for missing role")
	}
}

func TestOIDCLogin_Success(t *testing.T) {
	srv := newOIDCTestServer(t, http.StatusOK, "s.testtoken", []string{"default", "dev"})
	defer srv.Close()

	c, err := NewOIDCClient(srv.URL, "dev")
	if err != nil {
		t.Fatalf("NewOIDCClient: %v", err)
	}

	resp, err := c.Login("eyJhbGciOiJSUzI1NiJ9.test.sig")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if resp.Token != "s.testtoken" {
		t.Errorf("expected token s.testtoken, got %s", resp.Token)
	}
	if resp.LeaseDuration != 3600 {
		t.Errorf("expected lease 3600, got %d", resp.LeaseDuration)
	}
	if len(resp.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(resp.Policies))
	}
}

func TestOIDCLogin_Forbidden(t *testing.T) {
	srv := newOIDCTestServer(t, http.StatusForbidden, "", nil)
	defer srv.Close()

	c, err := NewOIDCClient(srv.URL, "dev")
	if err != nil {
		t.Fatalf("NewOIDCClient: %v", err)
	}

	_, err = c.Login("bad.jwt.token")
	if err == nil {
		t.Fatal("expected error for forbidden response")
	}
}

func TestOIDCLogin_EmptyJWT(t *testing.T) {
	c, err := NewOIDCClient("http://localhost:8200", "dev")
	if err != nil {
		t.Fatalf("NewOIDCClient: %v", err)
	}
	_, err = c.Login("")
	if err == nil {
		t.Fatal("expected error for empty jwt")
	}
}
