package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// PolicyRule represents a single ACL rule within a Vault policy.
type PolicyRule struct {
	Path         string   `json:"path"`
	Capabilities []string `json:"capabilities"`
}

// Policy represents a Vault ACL policy.
type Policy struct {
	Name  string       `json:"name"`
	Rules []PolicyRule `json:"rules"`
}

// policyResponse is the raw API response from Vault's policy endpoint.
type policyResponse struct {
	Data struct {
		Name   string `json:"name"`
		Policy string `json:"policy"`
	} `json:"data"`
}

// FetchPolicy retrieves the named policy from Vault and returns its metadata.
// It returns an error if the policy does not exist or the request fails.
func (c *Client) FetchPolicy(ctx context.Context, name string) (*Policy, error) {
	path := fmt.Sprintf("/v1/sys/policy/%s", name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.address+path, nil)
	if err != nil {
		return nil, fmt.Errorf("vault: build policy request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault: policy request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("vault: policy %q not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault: unexpected status %d fetching policy %q", resp.StatusCode, name)
	}

	var pr policyResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("vault: decode policy response: %w", err)
	}

	return &Policy{
		Name:  pr.Data.Name,
		Rules: []PolicyRule{},
	}, nil
}
