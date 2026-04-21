package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// LDAPClient authenticates against Vault using the LDAP auth method.
type LDAPClient struct {
	address  string
	mount    string
	httpClient *http.Client
}

// LDAPLoginResponse holds the token returned after a successful LDAP login.
type LDAPLoginResponse struct {
	Token    string
	Policies []string
}

// NewLDAPClient creates a new LDAPClient. mount defaults to "ldap" if empty.
func NewLDAPClient(address, mount string) (*LDAPClient, error) {
	if address == "" {
		return nil, fmt.Errorf("ldap: vault address is required")
	}
	if mount == "" {
		mount = "ldap"
	}
	return &LDAPClient{
		address:    strings.TrimRight(address, "/"),
		mount:      mount,
		httpClient: &http.Client{},
	}, nil
}

// Login authenticates the given username and password via Vault LDAP auth.
func (c *LDAPClient) Login(username, password string) (*LDAPLoginResponse, error) {
	if username == "" {
		return nil, fmt.Errorf("ldap: username is required")
	}
	if password == "" {
		return nil, fmt.Errorf("ldap: password is required")
	}

	url := fmt.Sprintf("%s/v1/auth/%s/login/%s", c.address, c.mount, username)
	body := fmt.Sprintf(`{"password":%q}`, password)

	resp, err := c.httpClient.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ldap: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("ldap: authentication failed (status %d)", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ldap: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string   `json:"client_token"`
			Policies    []string `json:"policies"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ldap: failed to decode response: %w", err)
	}

	return &LDAPLoginResponse{
		Token:    result.Auth.ClientToken,
		Policies: result.Auth.Policies,
	}, nil
}
