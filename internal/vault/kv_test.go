package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newKVTestServer(t *testing.T, version string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v1/sys/mounts/secret/tune":
			opts := map[string]interface{}{}
			if version != "" {
				opts["version"] = version
			}
			json.NewEncoder(w).Encode(map[string]interface{}{"data": nil, "options": opts})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestNewKVClient_EmptyMount(t *testing.T) {
	c := &Client{}
	_, err := NewKVClient(context.Background(), c, "")
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}

func TestNewKVClient_DefaultsToV1OnFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c, err := New(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	kv, err := NewKVClient(context.Background(), c, "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kv.Version() != KVv1 {
		t.Errorf("expected KVv1, got %v", kv.Version())
	}
}

func TestKVClient_ReadPath_V1(t *testing.T) {
	kv := &KVClient{version: KVv1, mount: "secret"}
	got := kv.ReadPath("myapp/db")
	want := "secret/myapp/db"
	if got != want {
		t.Errorf("ReadPath v1: got %q, want %q", got, want)
	}
}

func TestKVClient_ReadPath_V2(t *testing.T) {
	kv := &KVClient{version: KVv2, mount: "secret"}
	got := kv.ReadPath("myapp/db")
	want := "secret/data/myapp/db"
	if got != want {
		t.Errorf("ReadPath v2: got %q, want %q", got, want)
	}
}

func TestKVClient_MetaPath_V2(t *testing.T) {
	kv := &KVClient{version: KVv2, mount: "secret"}
	got := kv.MetaPath("myapp/db")
	want := "secret/metadata/myapp/db"
	if got != want {
		t.Errorf("MetaPath v2: got %q, want %q", got, want)
	}
}

func TestKVClient_MetaPath_V1(t *testing.T) {
	kv := &KVClient{version: KVv1, mount: "secret"}
	got := kv.MetaPath("myapp/db")
	want := "secret/myapp/db"
	if got != want {
		t.Errorf("MetaPath v1: got %q, want %q", got, want)
	}
}
