package vault

import (
	"sync"
	"time"
)

// CacheEntry holds a cached secret payload with an expiry timestamp.
type CacheEntry struct {
	Data      map[string]string
	FetchedAt time.Time
	TTL       time.Duration
}

// IsExpired reports whether the cache entry has passed its TTL.
func (e *CacheEntry) IsExpired() bool {
	return time.Since(e.FetchedAt) > e.TTL
}

// SecretCache is a thread-safe in-memory cache for Vault secret payloads.
type SecretCache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

// NewSecretCache creates a SecretCache with the given TTL for all entries.
func NewSecretCache(ttl time.Duration) *SecretCache {
	return &SecretCache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
}

// Get returns the cached data for the given path, or (nil, false) if absent or expired.
func (c *SecretCache) Get(path string) (map[string]string, bool) {
	c.mu.RLock()
	entry, ok := c.entries[path]
	c.mu.RUnlock()

	if !ok || entry.IsExpired() {
		return nil, false
	}
	return entry.Data, true
}

// Set stores secret data for the given path, resetting its TTL.
func (c *SecretCache) Set(path string, data map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[path] = &CacheEntry{
		Data:      data,
		FetchedAt: time.Now(),
		TTL:       c.ttl,
	}
}

// Invalidate removes the cached entry for the given path.
func (c *SecretCache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, path)
}

// Flush clears all cached entries.
func (c *SecretCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*CacheEntry)
}
