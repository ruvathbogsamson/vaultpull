// Package vault provides clients for interacting with HashiCorp Vault.
//
// # Transit Engine
//
// The TransitClient wraps Vault's Transit secrets engine, which provides
// encryption-as-a-service without exposing raw key material.
//
// Usage:
//
//	client, err := vault.NewTransitClient(address, token, "transit")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	cipher, err := client.Encrypt("my-key", "plaintext secret")
//	plain, err  := client.Decrypt("my-key", cipher)
//
// The mount parameter selects which Transit engine mount to target;
// it defaults to "transit" when left empty.
package vault
