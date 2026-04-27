package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newSentinelTestServer(t *testing.T, name string, policy *SentinelPolicy, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/sys/policies/egp/" + name
		if r.URL.Path != expected {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(status)
		if status == http.StatusOK && policy != nil {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": policy})
		}
	}))
}

func TestNewSentinelClient_MissingAddress(t *testing.T) {
	_, err := NewSentinelClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewSentinelClient_MissingToken(t *testing.T) {
	_, err := NewSentinelClient("http://localhost:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestGetPolicy_Success(t *testing.T) {
	policy := &SentinelPolicy{Name: "test-policy", Type: "egp", Body: "main = rule { true }"}
	srv := newSentinelTestServer(t, "test-policy", policy, http.StatusOK)
	defer srv.Close()

	client, err := NewSentinelClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := client.GetPolicy("test-policy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != policy.Name {
		t.Errorf("expected name %q, got %q", policy.Name, got.Name)
	}
	if got.Body != policy.Body {
		t.Errorf("expected body %q, got %q", policy.Body, got.Body)
	}
}

func TestGetPolicy_NotFound(t *testing.T) {
	srv := newSentinelTestServer(t, "missing", nil, http.StatusNotFound)
	defer srv.Close()

	client, err := NewSentinelClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.GetPolicy("missing")
	if err == nil {
		t.Fatal("expected error for not found policy")
	}
}

func TestGetPolicy_EmptyName(t *testing.T) {
	client, err := NewSentinelClient("http://localhost:8200", "token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = client.GetPolicy("")
	if err == nil {
		t.Fatal("expected error for empty policy name")
	}
}
