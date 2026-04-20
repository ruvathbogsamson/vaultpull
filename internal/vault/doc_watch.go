// Package vault provides a Vault API client with helpers for fetching,
// caching, diffing, renewing, and watching secrets.
//
// # Watcher
//
// Watcher polls a KV-v2 secret path at a configurable interval and fires
// callbacks when the secret data changes or an error occurs. It is safe to
// run multiple watchers concurrently for different paths.
//
// Basic usage:
//
//	client, _ := vault.New(addr, token)
//
//	w := vault.NewWatcher(client, "secret/data/myapp", vault.WatchConfig{
//		Interval: 60 * time.Second,
//		OnChange: func(path string, secrets map[string]string) {
//			fmt.Printf("secrets changed at %s\n", path)
//		},
//	})
//	w.Start(ctx)
//	defer w.Stop()
//
// The watcher compares the full secret map on each poll; any addition,
// removal, or value change triggers the OnChange callback.
package vault
