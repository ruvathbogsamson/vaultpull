package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newMetadataTestServer(status int, payload interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestFetchMetadata_Success(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"created_time":    now.Format(time.RFC3339),
			"updated_time":    now.Format(time.RFC3339),
			"current_version": 3,
			"oldest_version":  1,
			"custom_metadata": map[string]string{"owner": "team-a"},
		},
	}
	srv := newMetadataTestServer(http.StatusOK, payload)
	defer srv.Close()

	client := NewMetadataClient(srv.URL, "test-token", "secret")
	meta, err := client.FetchMetadata("myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.CurrentVersion != 3 {
		t.Errorf("expected current_version 3, got %d", meta.CurrentVersion)
	}
	if meta.OldestVersion != 1 {
		t.Errorf("expected oldest_version 1, got %d", meta.OldestVersion)
	}
	if meta.CustomMetadata["owner"] != "team-a" {
		t.Errorf("expected custom_metadata owner=team-a, got %q", meta.CustomMetadata["owner"])
	}
}

func TestFetchMetadata_NotFound(t *testing.T) {
	srv := newMetadataTestServer(http.StatusNotFound, nil)
	defer srv.Close()

	client := NewMetadataClient(srv.URL, "test-token", "secret")
	_, err := client.FetchMetadata("missing/path")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestFetchMetadata_ServerError(t *testing.T) {
	srv := newMetadataTestServer(http.StatusInternalServerError, nil)
	defer srv.Close()

	client := NewMetadataClient(srv.URL, "test-token", "secret")
	_, err := client.FetchMetadata("some/path")
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}

func TestFetchMetadata_Unreachable(t *testing.T) {
	client := NewMetadataClient("http://127.0.0.1:19999", "tok", "secret")
	_, err := client.FetchMetadata("any/path")
	if err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
