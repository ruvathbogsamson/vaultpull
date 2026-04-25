package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// PluginInfo holds metadata about a registered Vault plugin.
type PluginInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Builtin bool   `json:"builtin"`
	SHA256  string `json:"sha256"`
}

// PluginClient interacts with the Vault plugin catalog API.
type PluginClient struct {
	address string
	token   string
	httpClient *http.Client
}

// NewPluginClient creates a new PluginClient.
// Returns an error if address or token is empty.
func NewPluginClient(address, token string) (*PluginClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &PluginClient{
		address:    address,
		token:      token,
		httpClient: &http.Client{},
	}, nil
}

// ListPlugins returns all registered plugins of a given type (auth, secret, database).
// Pass an empty string to list all types.
func (c *PluginClient) ListPlugins(pluginType string) ([]PluginInfo, error) {
	path := "/v1/sys/plugins/catalog"
	if pluginType != "" {
		path = fmt.Sprintf("/v1/sys/plugins/catalog/%s", pluginType)
	}

	req, err := http.NewRequest(http.MethodGet, c.address+path, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("plugin type %q not found", pluginType)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Plugins []PluginInfo `json:"detailed"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return result.Data.Plugins, nil
}
