package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newAppRoleTestServer(t *testing.T, statusCode int, token string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/approle/login" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
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

func TestAppRoleLogin_Success(t *testing.T) {
	srv := newAppRoleTestServer(t, http.StatusOK, "s.testtoken")
	defer srv.Close()

	client := NewAppRoleClient(srv.URL, "")
	result, err := client.Login(context.Background(), AppRoleCredentials{
		RoleID:   "my-role",
		SecretID: "my-secret",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ClientToken != "s.testtoken" {
		t.Errorf("expected token %q, got %q", "s.testtoken", result.ClientToken)
	}
	if result.LeaseDuration != 3600 {
		t.Errorf("expected lease duration 3600, got %d", result.LeaseDuration)
	}
	if !result.Renewable {
		t.Error("expected renewable to be true")
	}
}

func TestAppRoleLogin_MissingRoleID(t *testing.T) {
	client := NewAppRoleClient("http://localhost", "")
	_, err := client.Login(context.Background(), AppRoleCredentials{SecretID: "s"})
	if err == nil {
		t.Fatal("expected error for missing role_id")
	}
}

func TestAppRoleLogin_MissingSecretID(t *testing.T) {
	client := NewAppRoleClient("http://localhost", "")
	_, err := client.Login(context.Background(), AppRoleCredentials{RoleID: "r"})
	if err == nil {
		t.Fatal("expected error for missing secret_id")
	}
}

func TestAppRoleLogin_ServerError(t *testing.T) {
	srv := newAppRoleTestServer(t, http.StatusForbidden, "")
	defer srv.Close()

	client := NewAppRoleClient(srv.URL, "approle")
	_, err := client.Login(context.Background(), AppRoleCredentials{
		RoleID:   "r",
		SecretID: "s",
	})
	if err == nil {
		t.Fatal("expected error on non-200 response")
	}
}

func TestNewAppRoleClient_DefaultMount(t *testing.T) {
	client := NewAppRoleClient("http://vault", "")
	if client.mountPath != "approle" {
		t.Errorf("expected default mount 'approle', got %q", client.mountPath)
	}
}

func TestNewAppRoleClient_CustomMount(t *testing.T) {
	client := NewAppRoleClient("http://vault", "custom-approle")
	if client.mountPath != "custom-approle" {
		t.Errorf("expected mount 'custom-approle', got %q", client.mountPath)
	}
}
