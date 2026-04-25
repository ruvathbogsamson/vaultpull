// Package vault provides clients for interacting with HashiCorp Vault.
//
// # Seal Client
//
// The SealClient provides operations for inspecting and managing the
// seal state of a Vault instance.
//
// # Usage
//
//	client, err := vault.NewSealClient(address, token)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Check current seal status
//	status, err := client.GetSealStatus()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sealed: %v\n", status.Sealed)
//
//	// Seal the vault (requires sudo policy)
//	if err := client.Seal(); err != nil {
//		log.Fatal(err)
//	}
package vault
