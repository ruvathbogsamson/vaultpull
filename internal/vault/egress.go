package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// EgressRule represents a single egress policy rule attached to a secret path.
type EgressRule struct {
	Path        string            `json:"path"`
	Capabilities []string         `json:"capabilities"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// EgressClient fetches egress rules from a Vault sys/egress endpoint.
type EgressClient struct {
	address string
	token   string
	httpClient *http.Client
}

// NewEgressClient creates a new EgressClient. Returns an error if address or token is empty.
func NewEgressClient(address, token string) (*EgressClient, error) {
	if strings.TrimSpace(address) == "" {
		return nil, fmt.Errorf("egress: vault address is required")
	}
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("egress: vault token is required")
	}
	return &EgressClient{
		address:    strings.TrimRight(address, "/"),
		token:      token,
		httpClient: &http.Client{},
	}, nil
}

// ListRules returns all egress rules from Vault for the given namespace.
func (c *EgressClient) ListRules(namespace string) ([]EgressRule, error) {
	url := fmt.Sprintf("%s/v1/sys/egress/%s", c.address, strings.Trim(namespace, "/"))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("egress: failed to build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("egress: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("egress: namespace %q not found", namespace)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("egress: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Rules []EgressRule `json:"rules"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("egress: failed to decode response: %w", err)
	}
	return result.Data.Rules, nil
}
