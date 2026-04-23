package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newPKITestServer(t *testing.T, role string, statusCode int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/pki/issue/" + role
		if r.URL.Path != expected {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if statusCode != http.StatusOK {
			http.Error(w, "error", statusCode)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"serial_number": "01:02:03",
				"certificate":   "-----BEGIN CERTIFICATE-----",
				"private_key":   "-----BEGIN RSA PRIVATE KEY-----",
				"issuing_ca":    "-----BEGIN CERTIFICATE-----CA",
				"expiration":    time.Now().Add(24 * time.Hour).Unix(),
			},
		})
	}))
}

func TestNewPKIClient_MissingAddress(t *testing.T) {
	_, err := NewPKIClient("", "token", "")
	if err == nil || !strings.Contains(err.Error(), "address") {
		t.Fatalf("expected address error, got %v", err)
	}
}

func TestNewPKIClient_MissingToken(t *testing.T) {
	_, err := NewPKIClient("http://localhost", "", "")
	if err == nil || !strings.Contains(err.Error(), "token") {
		t.Fatalf("expected token error, got %v", err)
	}
}

func TestNewPKIClient_DefaultMount(t *testing.T) {
	c, err := NewPKIClient("http://localhost", "tok", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.mount != "pki" {
		t.Errorf("expected default mount 'pki', got %q", c.mount)
	}
}

func TestIssueCertificate_Success(t *testing.T) {
	srv := newPKITestServer(t, "web", http.StatusOK)
	defer srv.Close()

	c, err := NewPKIClient(srv.URL, "test-token", "pki")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cert, err := c.IssueCertificate(IssueCertRequest{
		Role:       "web",
		CommonName: "example.com",
		TTL:        "24h",
		AltNames:   []string{"www.example.com"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cert.SerialNumber != "01:02:03" {
		t.Errorf("unexpected serial: %s", cert.SerialNumber)
	}
	if cert.Expiration.IsZero() {
		t.Error("expiration should not be zero")
	}
}

func TestIssueCertificate_MissingRole(t *testing.T) {
	c, _ := NewPKIClient("http://localhost", "tok", "pki")
	_, err := c.IssueCertificate(IssueCertRequest{CommonName: "example.com"})
	if err == nil || !strings.Contains(err.Error(), "role") {
		t.Fatalf("expected role error, got %v", err)
	}
}

func TestIssueCertificate_ServerError(t *testing.T) {
	srv := newPKITestServer(t, "web", http.StatusForbidden)
	defer srv.Close()

	c, _ := NewPKIClient(srv.URL, "tok", "pki")
	_, err := c.IssueCertificate(IssueCertRequest{Role: "web", CommonName: "example.com"})
	if err == nil || !strings.Contains(err.Error(), "403") {
		t.Fatalf("expected 403 error, got %v", err)
	}
}
