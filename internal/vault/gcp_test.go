package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newGCPTestServer(t *testing.T, roleset string, status int, creds GCPCredentials) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/gcp/token/" + roleset
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if status == http.StatusOK {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": creds})
		}
	}))
}

func TestNewGCPClient_MissingAddress(t *testing.T) {
	_, err := NewGCPClient("", "token", "")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewGCPClient_MissingToken(t *testing.T) {
	_, err := NewGCPClient("http://localhost", "", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewGCPClient_DefaultMount(t *testing.T) {
	c, err := NewGCPClient("http://localhost", "tok", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.mount != "gcp" {
		t.Errorf("expected default mount 'gcp', got %q", c.mount)
	}
}

func TestGenerateOAuthToken_Success(t *testing.T) {
	creds := GCPCredentials{
		Token:          "ya29.abc123",
		ExpireTime:     "2099-01-01T00:00:00Z",
		ServiceAccount: "svc@project.iam.gserviceaccount.com",
	}
	srv := newGCPTestServer(t, "my-roleset", http.StatusOK, creds)
	defer srv.Close()

	c, _ := NewGCPClient(srv.URL, "test-token", "")
	got, err := c.GenerateOAuthToken("my-roleset")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Token != creds.Token {
		t.Errorf("expected token %q, got %q", creds.Token, got.Token)
	}
	if got.ServiceAccount != creds.ServiceAccount {
		t.Errorf("expected service account %q, got %q", creds.ServiceAccount, got.ServiceAccount)
	}
}

func TestGenerateOAuthToken_NotFound(t *testing.T) {
	srv := newGCPTestServer(t, "other-roleset", http.StatusOK, GCPCredentials{})
	defer srv.Close()

	c, _ := NewGCPClient(srv.URL, "test-token", "")
	_, err := c.GenerateOAuthToken("missing-roleset")
	if err == nil {
		t.Fatal("expected error for missing roleset")
	}
}

func TestGenerateOAuthToken_MissingRoleset(t *testing.T) {
	c, _ := NewGCPClient("http://localhost", "tok", "")
	_, err := c.GenerateOAuthToken("")
	if err == nil {
		t.Fatal("expected error for empty roleset")
	}
}
