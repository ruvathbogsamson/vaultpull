package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AliasInfo represents an entity alias in Vault's identity store.
type AliasInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	MountAccessor string `json:"mount_accessor"`
	MountType    string `json:"mount_type"`
	EntityID     string `json:"canonical_id"`
}

// AliasClient interacts with Vault's identity alias API.
type AliasClient struct {
	address string
	token   string
	client  *http.Client
}

// NewAliasClient creates a new AliasClient.
func NewAliasClient(address, token string) (*AliasClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &AliasClient{
		address: address,
		token:   token,
		client:  &http.Client{},
	}, nil
}

// GetAlias retrieves an entity alias by ID.
func (c *AliasClient) GetAlias(id string) (*AliasInfo, error) {
	if id == "" {
		return nil, fmt.Errorf("alias id is required")
	}
	url := fmt.Sprintf("%s/v1/identity/entity-alias/id/%s", c.address, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("alias %q not found", id)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Data AliasInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result.Data, nil
}
