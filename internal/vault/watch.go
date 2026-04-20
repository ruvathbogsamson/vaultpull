package vault

import (
	"context"
	"log"
	"time"
)

// WatchConfig holds configuration for the secret watcher.
type WatchConfig struct {
	Interval time.Duration
	OnChange func(path string, secrets map[string]string)
	OnError  func(path string, err error)
}

// Watcher polls a Vault secret path at a configured interval and
// invokes callbacks when the secret changes or an error occurs.
type Watcher struct {
	client *Client
	path   string
	cfg    WatchConfig
	stop   chan struct{}
}

// NewWatcher creates a Watcher for the given client and secret path.
func NewWatcher(client *Client, path string, cfg WatchConfig) *Watcher {
	if cfg.Interval <= 0 {
		cfg.Interval = 30 * time.Second
	}
	if cfg.OnError == nil {
		cfg.OnError = func(path string, err error) {
			log.Printf("vault/watcher: error watching %s: %v", path, err)
		}
	}
	return &Watcher{
		client: client,
		path:   path,
		cfg:    cfg,
		stop:   make(chan struct{}),
	}
}

// Start begins polling in the background. Cancel ctx or call Stop to halt.
func (w *Watcher) Start(ctx context.Context) {
	go w.loop(ctx)
}

// Stop signals the watcher to cease polling.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) loop(ctx context.Context) {
	var prev map[string]string
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stop:
			return
		case <-ticker.C:
			current, err := FetchSecrets(ctx, w.client, w.path)
			if err != nil {
				w.cfg.OnError(w.path, err)
				continue
			}
			if hasChanged(prev, current) {
				w.cfg.OnChange(w.path, current)
				prev = current
			}
		}
	}
}

func hasChanged(prev, current map[string]string) bool {
	if len(prev) != len(current) {
		return true
	}
	for k, v := range current {
		if prev[k] != v {
			return true
		}
	}
	return false
}
