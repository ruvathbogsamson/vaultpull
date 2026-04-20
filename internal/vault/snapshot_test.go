package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newSnapshotTestServer(t *testing.T, path string, payload map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"data": payload}); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
}

func TestCapture_Success(t *testing.T) {
	payload := map[string]interface{}{"data": map[string]interface{}{"API_KEY": "abc123", "DB_PASS": "secret"}}
	srv := newSnapshotTestServer(t, "/v1/secret/data/myapp", payload)
	defer srv.Close()

	c, err := New(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	sc := NewSnapshotClient(c)
	snap, err := sc.Capture("secret/myapp")
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}

	if snap.Path != "secret/myapp" {
		t.Errorf("expected path %q, got %q", "secret/myapp", snap.Path)
	}
	if snap.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
	if snap.Data["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", snap.Data["API_KEY"])
	}
}

func TestSnapshot_MarshalRoundtrip(t *testing.T) {
	orig := &Snapshot{
		Path:       "secret/app",
		Data:       map[string]string{"FOO": "bar"},
		CapturedAt: time.Now().UTC().Truncate(time.Second),
	}

	b, err := orig.Marshal()
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	got, err := UnmarshalSnapshot(b)
	if err != nil {
		t.Fatalf("UnmarshalSnapshot: %v", err)
	}

	if got.Path != orig.Path {
		t.Errorf("path mismatch: want %q got %q", orig.Path, got.Path)
	}
	if got.Data["FOO"] != "bar" {
		t.Errorf("data mismatch: want bar got %q", got.Data["FOO"])
	}
	if !got.CapturedAt.Equal(orig.CapturedAt) {
		t.Errorf("timestamp mismatch: want %v got %v", orig.CapturedAt, got.CapturedAt)
	}
}

func TestFlattenSecretData_Nested(t *testing.T) {
	raw := map[string]interface{}{
		"data": map[string]interface{}{"KEY": "value"},
		"metadata": map[string]interface{}{"version": 2},
	}
	out, err := flattenSecretData(raw)
	if err != nil {
		t.Fatalf("flattenSecretData: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", out["KEY"])
	}
	if _, ok := out["metadata"]; ok {
		t.Error("metadata should not appear in flattened output")
	}
}
