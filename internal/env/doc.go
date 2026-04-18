// Package env provides utilities for filtering Vault secrets by namespace
// and writing them to .env files in KEY=VALUE format.
//
// Filter accepts a map of secrets and a prefix string, returning only the
// entries whose keys begin with the given prefix.
//
// NewWriter creates a Writer that targets the specified file path. Calling
// Write on the Writer will create or truncate the file and emit each
// key-value pair on its own line in KEY=VALUE format.
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
