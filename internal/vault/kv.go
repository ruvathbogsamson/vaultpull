// Package vault provides HashiCorp Vault integration for vaultpull.
package vault

import (
	"context"
	"fmt"
	"strings"
)

// KVVersion represents the KV secrets engine version.
type KVVersion int

const (
	KVv1 KVVersion = 1
	KVv2 KVVersion = 2
)

// KVClient wraps a Client with KV version-aware path handling.
type KVClient struct {
	client  *Client
	version KVVersion
	mount   string
}

// NewKVClient creates a KVClient by detecting the KV engine version at the
// given mount path. Falls back to v1 if detection fails.
func NewKVClient(ctx context.Context, c *Client, mount string) (*KVClient, error) {
	if mount == "" {
		return nil, fmt.Errorf("kv: mount path must not be empty")
	}
	mount = strings.Trim(mount, "/")

	version, err := detectKVVersion(ctx, c, mount)
	if err != nil {
		// Default to v1 on detection failure.
		version = KVv1
	}
	return &KVClient{client: c, version: version, mount: mount}, nil
}

// ReadPath returns the resolved API path for reading a secret.
func (k *KVClient) ReadPath(secretPath string) string {
	secretPath = strings.Trim(secretPath, "/")
	if k.version == KVv2 {
		return fmt.Sprintf("%s/data/%s", k.mount, secretPath)
	}
	return fmt.Sprintf("%s/%s", k.mount, secretPath)
}

// MetaPath returns the resolved API path for reading secret metadata (v2 only).
func (k *KVClient) MetaPath(secretPath string) string {
	secretPath = strings.Trim(secretPath, "/")
	if k.version == KVv2 {
		return fmt.Sprintf("%s/metadata/%s", k.mount, secretPath)
	}
	return fmt.Sprintf("%s/%s", k.mount, secretPath)
}

// Version returns the detected KV engine version.
func (k *KVClient) Version() KVVersion { return k.version }

// detectKVVersion queries the sys/mounts endpoint to determine the KV version.
func detectKVVersion(ctx context.Context, c *Client, mount string) (KVVersion, error) {
	path := fmt.Sprintf("sys/mounts/%s/tune", mount)
	secret, err := c.vc.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return KVv1, fmt.Errorf("kv: mount tune read failed: %w", err)
	}
	if secret == nil {
		return KVv1, fmt.Errorf("kv: no data returned for mount %q", mount)
	}
	opts, ok := secret.Data["options"].(map[string]interface{})
	if !ok {
		return KVv1, nil
	}
	if v, _ := opts["version"].(string); v == "2" {
		return KVv2, nil
	}
	return KVv1, nil
}
