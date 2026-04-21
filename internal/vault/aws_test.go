package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newAWSTestServer(t *testing.T, role string, statusCode int, creds *AWSCredentials) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/aws/creds/" + role
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(statusCode)
		if creds != nil {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data":           creds,
				"lease_duration": creds.LeaseDuration,
			})
		}
	}))
}

func TestNewAWSClient_MissingAddress(t *testing.T) {
	_, err := NewAWSClient("", "token", "")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewAWSClient_MissingToken(t *testing.T) {
	_, err := NewAWSClient("http://localhost", "", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewAWSClient_DefaultMount(t *testing.T) {
	c, err := NewAWSClient("http://localhost", "token", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.mount != "aws" {
		t.Errorf("expected default mount 'aws', got %q", c.mount)
	}
}

func TestGenerateCredentials_Success(t *testing.T) {
	want := &AWSCredentials{
		AccessKey:     "AKIAIOSFODNN7EXAMPLE",
		SecretKey:     "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		SecurityToken: "token123",
		LeaseDuration: 3600,
	}
	srv := newAWSTestServer(t, "my-role", http.StatusOK, want)
	defer srv.Close()

	c, _ := NewAWSClient(srv.URL, "test-token", "")
	got, err := c.GenerateCredentials("my-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.AccessKey != want.AccessKey {
		t.Errorf("access key: got %q, want %q", got.AccessKey, want.AccessKey)
	}
	if got.LeaseDuration != want.LeaseDuration {
		t.Errorf("lease duration: got %d, want %d", got.LeaseDuration, want.LeaseDuration)
	}
}

func TestGenerateCredentials_NotFound(t *testing.T) {
	srv := newAWSTestServer(t, "other-role", http.StatusNotFound, nil)
	defer srv.Close()

	c, _ := NewAWSClient(srv.URL, "test-token", "")
	_, err := c.GenerateCredentials("missing-role")
	if err == nil {
		t.Fatal("expected error for missing role")
	}
}

func TestGenerateCredentials_MissingRole(t *testing.T) {
	c, _ := NewAWSClient("http://localhost", "token", "")
	_, err := c.GenerateCredentials("")
	if err == nil {
		t.Fatal("expected error for empty role")
	}
}
