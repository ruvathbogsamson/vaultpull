// Package vault provides utilities for interacting with HashiCorp Vault.
//
// # Token Renewal
//
// TokenRenewer runs a background goroutine that periodically calls
// Vault's auth/token/renew-self endpoint so that long-running vaultpull
// processes do not lose access mid-sync.
//
// Usage:
//
//	renewer := vault.NewTokenRenewer(client, 5*time.Minute)
//	renewer.Start(ctx)
//	defer renewer.Stop()
//
// The renewer respects context cancellation, making it safe to use with
// the application's root context.
package vault
