// Package vault — KV engine version detection and path resolution.
//
// # KV Version Handling
//
// HashiCorp Vault supports two versions of the KV secrets engine:
//
//   - KV v1: secrets are stored at <mount>/<path>
//   - KV v2: secrets are stored at <mount>/data/<path>, with metadata
//     available at <mount>/metadata/<path>
//
// KVClient wraps a Vault client and automatically resolves read and
// metadata paths based on the detected engine version.
//
// # Version Detection
//
// Version detection queries the sys/mounts/<mount>/tune endpoint. If
// detection fails (e.g. insufficient permissions or network error), the
// client defaults to KV v1 to maintain backward compatibility.
//
// # Usage
//
//	kv, err := vault.NewKVClient(ctx, client, "secret")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	path := kv.ReadPath("myapp/credentials")
package vault
