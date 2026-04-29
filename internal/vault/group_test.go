package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newGroupTestServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestNewGroupClient_MissingAddress(t *testing.T) {
	_, err := NewGroupClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewGroupClient_MissingToken(t *testing.T) {
	_, err := NewGroupClient("http://localhost:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestGetGroup_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"id":   "abc-123",
			"name": "dev-team",
			"type": "internal",
			"policies": []string{"default", "dev"},
			"member_entity_ids": []string{"eid-1"},
			"metadata": map[string]string{"env": "staging"},
		},
	}
	srv := newGroupTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c, err := NewGroupClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	g, err := c.GetGroup("dev-team")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.ID != "abc-123" {
		t.Errorf("expected id abc-123, got %s", g.ID)
	}
	if g.Name != "dev-team" {
		t.Errorf("expected name dev-team, got %s", g.Name)
	}
	if len(g.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(g.Policies))
	}
}

func TestGetGroup_NotFound(t *testing.T) {
	srv := newGroupTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c, err := NewGroupClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = c.GetGroup("missing-group")
	if err == nil {
		t.Fatal("expected not-found error")
	}
}

func TestGetGroup_EmptyName(t *testing.T) {
	c, err := NewGroupClient("http://localhost:8200", "token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = c.GetGroup("")
	if err == nil {
		t.Fatal("expected error for empty group name")
	}
}
