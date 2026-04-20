package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newHealthTestServer(statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/sys/health" {
			w.WriteHeader(statusCode)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestCheckHealth_Initialized(t *testing.T) {
	srv := newHealthTestServer(http.StatusOK)
	defer srv.Close()

	c := &Client{address: srv.URL}
	status, err := c.CheckHealth(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Reachable || !status.Initialized || status.Sealed {
		t.Errorf("expected initialized unsealed vault, got %+v", status)
	}
}

func TestCheckHealth_Sealed(t *testing.T) {
	srv := newHealthTestServer(503)
	defer srv.Close()

	c := &Client{address: srv.URL}
	status, err := c.CheckHealth(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Sealed {
		t.Errorf("expected sealed vault, got %+v", status)
	}
}

func TestCheckHealth_Standby(t *testing.T) {
	srv := newHealthTestServer(http.StatusTooManyRequests)
	defer srv.Close()

	c := &Client{address: srv.URL}
	status, err := c.CheckHealth(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Standby {
		t.Errorf("expected standby vault, got %+v", status)
	}
}

func TestCheckHealth_Unreachable(t *testing.T) {
	c := &Client{address: "http://127.0.0.1:19999"}
	status, err := c.CheckHealth(context.Background())
	if err == nil {
		t.Fatal("expected error for unreachable vault")
	}
	if status.Reachable {
		t.Error("expected reachable=false")
	}
}
