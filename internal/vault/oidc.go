package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// OIDCClient handles OIDC authentication against a Vault server.
type OIDCClient struct {
	address string
	role    string
	client  *http.Client
}

// OIDCLoginResponse holds the result of a successful OIDC login.
type OIDCLoginResponse struct {
	Token    string
	Policies []string
	LeaseDuration int
}

// NewOIDCClient creates a new OIDCClient.
// address and role must be non-empty.
func NewOIDCClient(address, role string) (*OIDCClient, error) {
	address = strings.TrimRight(address, "/")
	if address == "" {
		return nil, fmt.Errorf("oidc: vault address is required")
	}
	if role == "" {
		return nil, fmt.Errorf("oidc: role is required")
	}
	return &OIDCClient{
		address: address,
		role:    role,
		client:  &http.Client{},
	}, nil
}

// Login exchanges a JWT token for a Vault client token via the OIDC auth method.
func (c *OIDCClient) Login(jwt string) (*OIDCLoginResponse, error) {
	if jwt == "" {
		return nil, fmt.Errorf("oidc: jwt is required")
	}

	payload := fmt.Sprintf(`{"jwt":%q,"role":%q}`, jwt, c.role)
	url := fmt.Sprintf("%s/v1/auth/oidc/login", c.address)

	resp, err := c.client.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("oidc: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("oidc: authentication failed (status %d)", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("oidc: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Auth struct {
			ClientToken   string   `json:"client_token"`
			Policies      []string `json:"policies"`
			LeaseDuration int      `json:"lease_duration"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("oidc: failed to decode response: %w", err)
	}

	return &OIDCLoginResponse{
		Token:         body.Auth.ClientToken,
		Policies:      body.Auth.Policies,
		LeaseDuration: body.Auth.LeaseDuration,
	}, nil
}
