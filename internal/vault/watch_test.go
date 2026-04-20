package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func newWatchTestServer(t *testing.T, responses []map[string]interface{}) *httptest.Server {
	t.Helper()
	var callCount atomic.Int32
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idx := int(callCount.Add(1)) - 1
		if idx >= len(responses) {
			idx = len(responses) - 1
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responses[idx])
	}))
}

func TestWatcher_DetectsChange(t *testing.T) {
	responses := []map[string]interface{}{
		{"data": map[string]interface{}{"data": map[string]interface{}{"KEY": "v1"}}},
		{"data": map[string]interface{}{"data": map[string]interface{}{"KEY": "v2"}}},
	}
	srv := newWatchTestServer(t, responses)
	defer srv.Close()

	client, err := New(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	changed := make(chan map[string]string, 2)
	cfg := WatchConfig{
		Interval: 20 * time.Millisecond,
		OnChange: func(_ string, s map[string]string) { changed <- s },
		OnError:  func(_ string, _ error) {},
	}

	w := NewWatcher(client, "secret/data/app", cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	w.Start(ctx)

	select {
	case s := <-changed:
		if s["KEY"] != "v1" && s["KEY"] != "v2" {
			t.Errorf("unexpected secret value: %v", s)
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for change callback")
	}
}

func TestWatcher_Stop(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"data": map[string]interface{}{"K": "v"}},
		})
	}))
	defer srv.Close()

	client, _ := New(srv.URL, "tok")
	cfg := WatchConfig{
		Interval: 20 * time.Millisecond,
		OnChange: func(_ string, _ map[string]string) {},
	}
	w := NewWatcher(client, "secret/data/app", cfg)
	w.Start(context.Background())
	time.Sleep(60 * time.Millisecond)
	w.Stop()
	snap := calls.Load()
	time.Sleep(60 * time.Millisecond)
	if calls.Load() > snap+1 {
		t.Errorf("watcher continued polling after Stop")
	}
}

func TestHasChanged(t *testing.T) {
	if hasChanged(nil, map[string]string{"A": "1"}) != true {
		t.Error("expected change when prev is nil")
	}
	if hasChanged(map[string]string{"A": "1"}, map[string]string{"A": "1"}) != false {
		t.Error("expected no change for identical maps")
	}
	if hasChanged(map[string]string{"A": "1"}, map[string]string{"A": "2"}) != true {
		t.Error("expected change when value differs")
	}
}
