package vault

import (
	"testing"
	"time"
)

func TestSecretCache_SetAndGet(t *testing.T) {
	c := NewSecretCache(5 * time.Minute)
	data := map[string]string{"KEY": "value"}

	c.Set("secret/app", data)

	got, ok := c.Get("secret/app")
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if got["KEY"] != "value" {
		t.Errorf("expected \"value\", got %q", got["KEY"])
	}
}

func TestSecretCache_Miss(t *testing.T) {
	c := NewSecretCache(5 * time.Minute)

	_, ok := c.Get("secret/missing")
	if ok {
		t.Fatal("expected cache miss, got hit")
	}
}

func TestSecretCache_Expiry(t *testing.T) {
	c := NewSecretCache(10 * time.Millisecond)
	c.Set("secret/app", map[string]string{"K": "v"})

	time.Sleep(20 * time.Millisecond)

	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected expired cache miss, got hit")
	}
}

func TestSecretCache_Invalidate(t *testing.T) {
	c := NewSecretCache(5 * time.Minute)
	c.Set("secret/app", map[string]string{"K": "v"})

	c.Invalidate("secret/app")

	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected miss after invalidation, got hit")
	}
}

func TestSecretCache_Flush(t *testing.T) {
	c := NewSecretCache(5 * time.Minute)
	c.Set("secret/a", map[string]string{"A": "1"})
	c.Set("secret/b", map[string]string{"B": "2"})

	c.Flush()

	for _, path := range []string{"secret/a", "secret/b"} {
		if _, ok := c.Get(path); ok {
			t.Errorf("expected miss for %q after flush, got hit", path)
		}
	}
}

func TestCacheEntry_IsExpired(t *testing.T) {
	entry := &CacheEntry{
		FetchedAt: time.Now().Add(-10 * time.Second),
		TTL:       5 * time.Second,
	}
	if !entry.IsExpired() {
		t.Error("expected entry to be expired")
	}

	fresh := &CacheEntry{
		FetchedAt: time.Now(),
		TTL:       5 * time.Minute,
	}
	if fresh.IsExpired() {
		t.Error("expected fresh entry to not be expired")
	}
}
