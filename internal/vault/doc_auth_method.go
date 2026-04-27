// Package vault provides a client for interacting with HashiCorp Vault.
//
// # Auth Method Client
//
// The AuthMethodClient allows callers to list all authentication methods
// currently enabled on a Vault server.
//
// Usage:
//
//	c, err := vault.NewAuthMethodClient(address, token)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	methods, err := c.ListAuthMethods()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for path, m := range methods {
//		fmt.Printf("%s -> %s\n", path, m.Type)
//	}
package vault
