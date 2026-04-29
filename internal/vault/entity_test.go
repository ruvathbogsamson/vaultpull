package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newEntityTestServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestNewEntityClient_MissingAddress(t *testing.T) {
	_, err := NewEntityClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewEntityClient_MissingToken(t *testing.T) {
	_, err := NewEntityClient("http://localhost:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestGetEntity_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"id":       "abc-123",
			"name":     "alice",
			"policies": []string{"default"},
			"metadata": map[string]string{"team": "eng"},
			"disabled": false,
		},
	}
	srv := newEntityTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c, err := NewEntityClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entity, err := c.GetEntity("alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entity.ID != "abc-123" {
		t.Errorf("expected id abc-123, got %s", entity.ID)
	}
	if entity.Name != "alice" {
		t.Errorf("expected name alice, got %s", entity.Name)
	}
}

func TestGetEntity_NotFound(t *testing.T) {
	srv := newEntityTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c, err := NewEntityClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = c.GetEntity("ghost")
	if err == nil {
		t.Fatal("expected not found error")
	}
}

func TestGetEntity_ServerError(t *testing.T) {
	srv := newEntityTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c, err := NewEntityClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = c.GetEntity("alice")
	if err == nil {
		t.Fatal("expected error for server error response")
	}
}
