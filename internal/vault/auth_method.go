package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AuthMethod represents a single auth method enabled in Vault.
type AuthMethod struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Accessor    string `json:"accessor"`
	Local       bool   `json:"local"`
}

// AuthMethodClient fetches enabled auth methods from Vault.
type AuthMethodClient struct {
	address string
	token   string
	client  *http.Client
}

// NewAuthMethodClient creates a new AuthMethodClient.
func NewAuthMethodClient(address, token string) (*AuthMethodClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &AuthMethodClient{
		address: address,
		token:   token,
		client:  &http.Client{},
	}, nil
}

// ListAuthMethods returns all enabled auth methods keyed by their mount path.
func (c *AuthMethodClient) ListAuthMethods() (map[string]AuthMethod, error) {
	url := fmt.Sprintf("%s/v1/sys/auth", c.address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("permission denied")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var result map[string]AuthMethod
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return result, nil
}
