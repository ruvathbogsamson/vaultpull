package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newKubernetesTestServer(t *testing.T, status int, token string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{
				"client_token": token,
				"accessor":     "test-accessor",
				"policies":     []string{"default", "k8s-policy"},
			},
		})
	}))
}

func TestNewKubernetesClient_MissingAddress(t *testing.T) {
	_, err := NewKubernetesClient("", "")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewKubernetesClient_DefaultMount(t *testing.T) {
	c, err := NewKubernetesClient("http://localhost:8200", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.mount != "kubernetes" {
		t.Errorf("expected default mount 'kubernetes', got %q", c.mount)
	}
}

func TestKubernetesLogin_Success(t *testing.T) {
	srv := newKubernetesTestServer(t, http.StatusOK, "k8s-token-abc")
	defer srv.Close()

	c, err := NewKubernetesClient(srv.URL, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := c.Login("my-role", "jwt-token-value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ClientToken != "k8s-token-abc" {
		t.Errorf("expected token 'k8s-token-abc', got %q", resp.ClientToken)
	}
	if len(resp.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(resp.Policies))
	}
}

func TestKubernetesLogin_MissingRole(t *testing.T) {
	c, _ := NewKubernetesClient("http://localhost:8200", "")
	_, err := c.Login("", "some-jwt")
	if err == nil {
		t.Fatal("expected error for missing role")
	}
}

func TestKubernetesLogin_MissingJWT(t *testing.T) {
	c, _ := NewKubernetesClient("http://localhost:8200", "")
	_, err := c.Login("my-role", "")
	if err == nil {
		t.Fatal("expected error for missing jwt")
	}
}

func TestKubernetesLogin_Forbidden(t *testing.T) {
	srv := newKubernetesTestServer(t, http.StatusForbidden, "")
	defer srv.Close()

	c, _ := NewKubernetesClient(srv.URL, "")
	_, err := c.Login("my-role", "bad-jwt")
	if err == nil {
		t.Fatal("expected error for forbidden response")
	}
}
