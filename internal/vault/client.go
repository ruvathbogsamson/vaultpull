package vault

import (
	"context"
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with namespace filtering.
type Client struct {
	api       *vaultapi.Client
	namespace string
}

// New creates a new Vault client using the provided address and token.
func New(address, token, namespace string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	api.SetToken(token)

	if namespace != "" {
		api.SetNamespace(namespace)
	}

	return &Client{api: api, namespace: namespace}, nil
}

// ReadSecrets reads key/value secrets from the given KV v2 path.
// It returns a map of secret key -> value strings.
func (c *Client) ReadSecrets(ctx context.Context, secretPath string) (map[string]string, error) {
	// Normalize path: strip leading slash
	secretPath = strings.TrimPrefix(secretPath, "/")

	secret, err := c.api.KVv2(mountFromPath(secretPath)).Get(ctx, dataPathFromPath(secretPath))
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", secretPath, err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found at path %q", secretPath)
	}

	result := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		str, ok := v.(string)
		if !ok {
			str = fmt.Sprintf("%v", v)
		}
		result[k] = str
	}

	return result, nil
}

// mountFromPath returns the KV mount (first path segment).
func mountFromPath(p string) string {
	parts := strings.SplitN(p, "/", 2)
	return parts[0]
}

// dataPathFromPath returns everything after the first path segment.
func dataPathFromPath(p string) string {
	parts := strings.SplitN(p, "/", 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}
