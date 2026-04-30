package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newAliasTestServer(status int, payload interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestNewAliasClient_MissingAddress(t *testing.T) {
	_, err := NewAliasClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewAliasClient_MissingToken(t *testing.T) {
	_, err := NewAliasClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestGetAlias_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"id":             "alias-123",
			"name":           "testuser",
			"mount_accessor": "auth_userpass_abc",
			"mount_type":     "userpass",
			"canonical_id":   "entity-456",
		},
	}
	srv := newAliasTestServer(http.StatusOK, payload)
	defer srv.Close()

	c, err := NewAliasClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	alias, err := c.GetAlias("alias-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alias.ID != "alias-123" {
		t.Errorf("expected id alias-123, got %s", alias.ID)
	}
	if alias.Name != "testuser" {
		t.Errorf("expected name testuser, got %s", alias.Name)
	}
	if alias.EntityID != "entity-456" {
		t.Errorf("expected entity entity-456, got %s", alias.EntityID)
	}
}

func TestGetAlias_NotFound(t *testing.T) {
	srv := newAliasTestServer(http.StatusNotFound, nil)
	defer srv.Close()

	c, _ := NewAliasClient(srv.URL, "test-token")
	_, err := c.GetAlias("missing")
	if err == nil {
		t.Fatal("expected error for not found alias")
	}
}

func TestGetAlias_ServerError(t *testing.T) {
	srv := newAliasTestServer(http.StatusInternalServerError, nil)
	defer srv.Close()

	c, _ := NewAliasClient(srv.URL, "test-token")
	_, err := c.GetAlias("alias-123")
	if err == nil {
		t.Fatal("expected error for server error")
	}
}

func TestGetAlias_EmptyID(t *testing.T) {
	c, _ := NewAliasClient("http://127.0.0.1:8200", "token")
	_, err := c.GetAlias("")
	if err == nil {
		t.Fatal("expected error for empty id")
	}
}
