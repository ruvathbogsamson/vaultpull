package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// TokenRenewer manages periodic renewal of a Vault token.
type TokenRenewer struct {
	client   *vaultapi.Client
	interval time.Duration
	stopCh   chan struct{}
}

// NewTokenRenewer creates a TokenRenewer that renews the token at the given interval.
// If interval is zero, a default of 5 minutes is used.
func NewTokenRenewer(client *vaultapi.Client, interval time.Duration) *TokenRenewer {
	if interval <= 0 {
		interval = 5 * time.Minute
	}
	return &TokenRenewer{
		client:   client,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start begins the renewal loop in a background goroutine.
// The loop stops when ctx is cancelled or Stop is called.
func (r *TokenRenewer) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := r.renew(); err != nil {
					// non-fatal: log and continue
					_ = fmt.Errorf("vault token renewal failed: %w", err)
				}
			case <-r.stopCh:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Stop halts the renewal loop.
func (r *TokenRenewer) Stop() {
	close(r.stopCh)
}

// renew calls the Vault token self-renew endpoint.
func (r *TokenRenewer) renew() error {
	_, err := r.client.Auth().Token().RenewSelf(0)
	if err != nil {
		return fmt.Errorf("renew-self: %w", err)
	}
	return nil
}
