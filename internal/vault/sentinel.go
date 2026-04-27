package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SentinelPolicy represents a Vault Sentinel policy entry.
type SentinelPolicy struct {
	Name enforcement string `json:"name"`
	Type string `json:"type"`
	Body string `json:"body"`
}

// SentinelClient interacts with Vault's Sentinel policy endpoints.
type SentinelClient struct {
	address string
	token   string
	httpClient *http.Client
}

// NewSentinelClient creates a new SentinelClient.
func NewSentinelClient(address, token string) (*SentinelClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &SentinelClient{
		address:    address,
		token:      token,
		httpClient: &http.Client{},
	}, nil
}

// GetPolicy retrieves a Sentinel policy by name.
func (c *SentinelClient) GetPolicy(name string) (*SentinelPolicy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name is required")
	}
	url := fmt.Sprintf("%s/v1/sys/policies/egp/%s", c.address, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("policy %q not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Data SentinelPolicy `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result.Data, nil
}
