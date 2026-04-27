package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newUnsealTestServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	}))
}

func TestNewUnsealClient_MissingAddress(t *testing.T) {
	_, err := NewUnsealClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewUnsealClient_MissingToken(t *testing.T) {
	_, err := NewUnsealClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestSubmitKey_Success(t *testing.T) {
	expected := UnsealResponse{Sealed: false, T: 3, N: 5, Progress: 3}
	srv := newUnsealTestServer(t, http.StatusOK, expected)
	defer srv.Close()

	c, err := NewUnsealClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewUnsealClient: %v", err)
	}

	resp, err := c.SubmitKey("abc123")
	if err != nil {
		t.Fatalf("SubmitKey: %v", err)
	}
	if resp.Sealed != false {
		t.Errorf("expected sealed=false, got %v", resp.Sealed)
	}
	if resp.Progress != 3 {
		t.Errorf("expected progress=3, got %d", resp.Progress)
	}
}

func TestSubmitKey_ServerError(t *testing.T) {
	srv := newUnsealTestServer(t, http.StatusInternalServerError, map[string]string{"errors": "internal"})
	defer srv.Close()

	c, err := NewUnsealClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewUnsealClient: %v", err)
	}

	_, err = c.SubmitKey("badkey")
	if err == nil {
		t.Fatal("expected error on server error")
	}
}

func TestReset_Success(t *testing.T) {
	expected := UnsealResponse{Sealed: true, Progress: 0}
	srv := newUnsealTestServer(t, http.StatusOK, expected)
	defer srv.Close()

	c, err := NewUnsealClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewUnsealClient: %v", err)
	}

	resp, err := c.Reset()
	if err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if resp.Progress != 0 {
		t.Errorf("expected progress=0 after reset, got %d", resp.Progress)
	}
}
