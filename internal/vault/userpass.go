package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// UserpassClient authenticates with Vault using the userpass auth method.
type UserpassClient struct {
	address string
	mount   string
	httpClient *http.Client
}

// UserpassCredentials holds the username and password for authentication.
type UserpassCredentials struct {
	Username string
	Password string
}

// UserpassToken is returned after a successful userpass login.
type UserpassToken struct {
	ClientToken string
	LeaseDuration int
	Renewable     bool
}

// NewUserpassClient creates a new UserpassClient.
// mount defaults to "userpass" if empty.
func NewUserpassClient(address, mount string) (*UserpassClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if mount == "" {
		mount = "userpass"
	}
	return &UserpassClient{
		address:    strings.TrimRight(address, "/"),
		mount:      mount,
		httpClient: &http.Client{},
	}, nil
}

// Login authenticates with Vault using username/password and returns a token.
func (c *UserpassClient) Login(creds UserpassCredentials) (*UserpassToken, error) {
	if creds.Username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if creds.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	url := fmt.Sprintf("%s/v1/auth/%s/login/%s", c.address, c.mount, creds.Username)
	body := fmt.Sprintf(`{"password":%q}`, creds.Password)

	resp, err := c.httpClient.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("userpass login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("userpass login failed: invalid credentials (status %d)", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userpass login failed: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken   string `json:"client_token"`
			LeaseDuration int    `json:"lease_duration"`
			Renewable     bool   `json:"renewable"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode userpass response: %w", err)
	}

	return &UserpassToken{
		ClientToken:   result.Auth.ClientToken,
		LeaseDuration: result.Auth.LeaseDuration,
		Renewable:     result.Auth.Renewable,
	}, nil
}
