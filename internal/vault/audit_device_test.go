package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newAuditDeviceTestServer(status int, body interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestNewAuditDeviceClient_MissingAddress(t *testing.T) {
	_, err := NewAuditDeviceClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewAuditDeviceClient_MissingToken(t *testing.T) {
	_, err := NewAuditDeviceClient("http://localhost:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestListAuditDevices_Success(t *testing.T) {
	payload := map[string]*AuditDevice{
		"file/": {
			Type:        "file",
			Description: "file audit log",
			Options:     map[string]string{"file_path": "/var/log/vault/audit.log"},
			Path:        "file/",
		},
	}
	srv := newAuditDeviceTestServer(http.StatusOK, payload)
	defer srv.Close()

	c, err := NewAuditDeviceClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	devices, err := c.ListAuditDevices()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := devices["file/"]; !ok {
		t.Error("expected file/ audit device in response")
	}
}

func TestListAuditDevices_Forbidden(t *testing.T) {
	srv := newAuditDeviceTestServer(http.StatusForbidden, nil)
	defer srv.Close()

	c, _ := NewAuditDeviceClient(srv.URL, "bad-token")
	_, err := c.ListAuditDevices()
	if err == nil {
		t.Fatal("expected error for forbidden response")
	}
}

func TestListAuditDevices_NotFound(t *testing.T) {
	srv := newAuditDeviceTestServer(http.StatusNotFound, nil)
	defer srv.Close()

	c, _ := NewAuditDeviceClient(srv.URL, "token")
	_, err := c.ListAuditDevices()
	if err == nil {
		t.Fatal("expected error for not found response")
	}
}
