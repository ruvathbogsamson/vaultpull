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
package sync
