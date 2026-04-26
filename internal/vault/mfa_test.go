package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMFATestServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestNewMFAClient_MissingAddress(t *testing.T) {
	_, err := NewMFAClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewMFAClient_MissingToken(t *testing.T) {
	_, err := NewMFAClient("http://localhost:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestMFAValidate_Success(t *testing.T) {
	payload := map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": "s.newtoken",
			"policies":     []string{"default"},
		},
	}
	ts := newMFATestServer(t, http.StatusOK, payload)
	defer ts.Close()

	c, err := NewMFAClient(ts.URL, "root")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := c.Validate(MFAValidateRequest{
		MFARequestID: "req-abc-123",
		MFAPayload:   map[string]string{"totp": "123456"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "s.newtoken" {
		t.Errorf("expected token s.newtoken, got %q", resp.Token)
	}
	if len(resp.Policies) != 1 || resp.Policies[0] != "default" {
		t.Errorf("unexpected policies: %v", resp.Policies)
	}
}

func TestMFAValidate_Forbidden(t *testing.T) {
	ts := newMFATestServer(t, http.StatusForbidden, nil)
	defer ts.Close()

	c, _ := NewMFAClient(ts.URL, "root")
	_, err := c.Validate(MFAValidateRequest{MFARequestID: "req-xyz"})
	if err == nil {
		t.Fatal("expected error for forbidden response")
	}
}

func TestMFAValidate_MissingRequestID(t *testing.T) {
	ts := newMFATestServer(t, http.StatusOK, nil)
	defer ts.Close()

	c, _ := NewMFAClient(ts.URL, "root")
	_, err := c.Validate(MFAValidateRequest{})
	if err == nil {
		t.Fatal("expected error for missing mfa_request_id")
	}
}
