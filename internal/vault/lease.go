package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// LeaseInfo holds metadata about a Vault secret lease.
type LeaseInfo struct {
	LeaseID   string
	Duration  time.Duration
	Renewable bool
	IssuedAt  time.Time
}

// LeaseManager tracks and renews secret leases obtained from Vault.
type LeaseManager struct {
	client *vaultapi.Client
	leases map[string]*LeaseInfo
}

// NewLeaseManager creates a LeaseManager backed by the given Vault client.
func NewLeaseManager(client *vaultapi.Client) *LeaseManager {
	return &LeaseManager{
		client: client,
		leases: make(map[string]*LeaseInfo),
	}
}

// Track registers a lease for management.
func (lm *LeaseManager) Track(leaseID string, duration time.Duration, renewable bool) {
	lm.leases[leaseID] = &LeaseInfo{
		LeaseID:   leaseID,
		Duration:  duration,
		Renewable: renewable,
		IssuedAt:  time.Now(),
	}
}

// Renew attempts to renew the lease with the given ID.
func (lm *LeaseManager) Renew(ctx context.Context, leaseID string, increment time.Duration) error {
	info, ok := lm.leases[leaseID]
	if !ok {
		return fmt.Errorf("lease %q not tracked", leaseID)
	}
	if !info.Renewable {
		return fmt.Errorf("lease %q is not renewable", leaseID)
	}

	seconds := int(increment.Seconds())
	secret, err := lm.client.Sys().RenewWithContext(ctx, leaseID, seconds)
	if err != nil {
		return fmt.Errorf("renew lease %q: %w", leaseID, err)
	}

	info.Duration = time.Duration(secret.LeaseDuration) * time.Second
	info.IssuedAt = time.Now()
	return nil
}

// Revoke revokes the lease and removes it from tracking.
func (lm *LeaseManager) Revoke(ctx context.Context, leaseID string) error {
	if err := lm.client.Sys().RevokeWithContext(ctx, leaseID); err != nil {
		return fmt.Errorf("revoke lease %q: %w", leaseID, err)
	}
	delete(lm.leases, leaseID)
	return nil
}

// Get returns tracked lease info, or nil if not found.
func (lm *LeaseManager) Get(leaseID string) *LeaseInfo {
	return lm.leases[leaseID]
}

// Expiring returns leases that expire within the given threshold.
func (lm *LeaseManager) Expiring(threshold time.Duration) []*LeaseInfo {
	var result []*LeaseInfo
	now := time.Now()
	for _, info := range lm.leases {
		expiry := info.IssuedAt.Add(info.Duration)
		if expiry.Sub(now) <= threshold {
			result = append(result, info)
		}
	}
	return result
}
