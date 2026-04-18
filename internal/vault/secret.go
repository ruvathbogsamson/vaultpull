package vault

import (
	"fmt"
	"strings"
)

// Secret represents a key-value secret fetched from Vault.
type Secret struct {
	Path string
	Data map[string]string
}

// FetchSecrets reads KV v2 secrets from the given path and returns a Secret.
func (c *Client) FetchSecrets(path string) (*Secret, error) {
	mount := mountFromPath(path)
	dataPath := dataPathFromPath(path)

	secret, err := c.logical.Read(fmt.Sprintf("%s/data/%s", mount, dataPath))
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", path)
	}

	rawData, ok := secret.Data["data"]
	if !ok {
		return nil, fmt.Errorf("secret at %q missing 'data' key", path)
	}

	rawMap, ok := rawData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format at %q", path)
	}

	result := &Secret{
		Path: path,
		Data: make(map[string]string, len(rawMap)),
	}

	for k, v := range rawMap {
		result.Data[strings.ToUpper(k)] = fmt.Sprintf("%v", v)
	}

	return result, nil
}
