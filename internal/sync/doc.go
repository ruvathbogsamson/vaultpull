// Package sync provides the top-level orchestration for vaultpull.
//
// It wires together the Vault client, namespace filtering, and .env file
// writing into a single [Syncer] type. Typical usage:
//
//	cfg, err := config.Load()
//	if err != nil { ... }
//
//	s, err := sync.New(cfg)
//	if err != nil { ... }
//
//	count, err := s.Run()
//	if err != nil { ... }
//	fmt.Printf("wrote %d secrets\n", count)
//
// # Error Handling
//
// Errors from individual secret paths are collected and returned as a combined
// error after all paths have been attempted, so a single unreachable path does
// not abort the entire sync. A fatal configuration or authentication error will
// still return immediately.
package sync
