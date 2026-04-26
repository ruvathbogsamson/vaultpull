package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newStepDownTestServer(statusCode int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("X-Vault-Token") == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(statusCode)
		if body != "" {
			_, _ = w.Write([]byte(body))
		}
	}))
}

func TestNewStepDownClient_MissingAddress(t *testing.T) {
	_, err := NewStepDownClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewStepDownClient_MissingToken(t *testing.T) {
	_, err := NewStepDownClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestStepDown_Success(t *testing.T) {
	srv := newStepDownTestServer(http.StatusNoContent, "")
	defer srv.Close()

	client, err := NewStepDownClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	status, err := client.StepDown()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Success {
		t.Errorf("expected success=true, got false")
	}
}

func TestStepDown_Forbidden(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"errors":["permission denied"]}`))
	}))
	defer srv.Close()

	client, err := NewStepDownClient(srv.URL, "bad-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.StepDown()
	if err == nil {
		t.Fatal("expected error for forbidden step-down")
	}
}

func TestStepDown_Unreachable(t *testing.T) {
	client, err := NewStepDownClient("http://127.0.0.1:19999", "token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = client.StepDown()
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
