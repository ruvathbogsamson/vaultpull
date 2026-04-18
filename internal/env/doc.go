// Package env provides utilities for filtering Vault secrets by namespace
// and writing them to .env files in KEY=VALUE format.
//
// Usage:
//
//	secrets := map[string]string{"APP_TOKEN": "abc", "DB_PASS": "xyz"}
//	filtered := env.Filter(secrets, "APP")
//	w := env.NewWriter(".env")
//	if err := w.Write(filtered); err != nil {
//		log.Fatal(err)
//	}
package env
