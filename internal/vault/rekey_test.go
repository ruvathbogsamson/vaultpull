package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newRekeyTestServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/rekey/init" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestNewRekeyClient_MissingAddress(t *testing.T) {
	_, err := NewRekeyClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewRekeyClient_MissingToken(t *testing.T) {
	_, err := NewRekeyClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestGetRekeyStatus_NotStarted(t *testing.T) {
	payload := RekeyStatus{Started: false, T: 0, N: 0}
	srv := newRekeyTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c, err := NewRekeyClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	status, err := c.GetRekeyStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Started {
		t.Error("expected rekey not started")
	}
}

func TestGetRekeyStatus_InProgress(t *testing.T) {
	payload := RekeyStatus{Started: true, T: 3, N: 5, Progress: 1, Required: 3, Nonce: "abc-123"}
	srv := newRekeyTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c, _ := NewRekeyClient(srv.URL, "test-token")
	status, err := c.GetRekeyStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Started {
		t.Error("expected rekey in progress")
	}
	if status.Nonce != "abc-123" {
		t.Errorf("expected nonce abc-123, got %s", status.Nonce)
	}
	if status.Progress != 1 {
		t.Errorf("expected progress 1, got %d", status.Progress)
	}
}

func TestGetRekeyStatus_ServerError(t *testing.T) {
	srv := newRekeyTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c, _ := NewRekeyClient(srv.URL, "test-token")
	_, err := c.GetRekeyStatus()
	if err == nil {
		t.Fatal("expected error on server error")
	}
}
