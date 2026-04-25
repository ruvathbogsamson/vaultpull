package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newQuotaTestServer(t *testing.T, name string, info *QuotaInfo, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/sys/quotas/rate-limit/" + name
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if info != nil {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": info})
		}
	}))
}

func TestNewQuotaClient_MissingAddress(t *testing.T) {
	_, err := NewQuotaClient("", "token")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewQuotaClient_MissingToken(t *testing.T) {
	_, err := NewQuotaClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestGetQuota_Success(t *testing.T) {
	info := &QuotaInfo{
		Name:        "global-rate",
		Path:        "secret/",
		Type:        "rate-limit",
		MaxRequests: 100,
		Rate:        50.0,
		Interval:    1.0,
	}
	srv := newQuotaTestServer(t, "global-rate", info, http.StatusOK)
	defer srv.Close()

	client, err := NewQuotaClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := client.GetQuota("global-rate")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != info.Name {
		t.Errorf("expected name %q, got %q", info.Name, got.Name)
	}
	if got.Rate != info.Rate {
		t.Errorf("expected rate %v, got %v", info.Rate, got.Rate)
	}
}

func TestGetQuota_NotFound(t *testing.T) {
	srv := newQuotaTestServer(t, "missing", nil, http.StatusNotFound)
	defer srv.Close()

	client, _ := NewQuotaClient(srv.URL, "test-token")
	_, err := client.GetQuota("missing")
	if err == nil {
		t.Fatal("expected not-found error")
	}
}

func TestGetQuota_EmptyName(t *testing.T) {
	client, _ := NewQuotaClient("http://127.0.0.1:8200", "token")
	_, err := client.GetQuota("")
	if err == nil {
		t.Fatal("expected error for empty quota name")
	}
}

func TestGetQuota_ServerError(t *testing.T) {
	srv := newQuotaTestServer(t, "bad", nil, http.StatusInternalServerError)
	defer srv.Close()

	client, _ := NewQuotaClient(srv.URL, "token")
	_, err := client.GetQuota("bad")
	if err == nil {
		t.Fatal("expected error on 500")
	}
}
