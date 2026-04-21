package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TokenInfo holds metadata about the current Vault token.
type TokenInfo struct {
	Accessor   string        `json:"accessor"`
	Policies   []string      `json:"policies"`
	TTL        time.Duration `json:"-"`
	Renewable  bool          `json:"renewable"`
	ExpireTime string        `json:"expire_time"`
}

// AuthTokenClient validates and inspects Vault tokens.
type AuthTokenClient struct {
	address string
	token   string
	client  *http.Client
}

// NewAuthTokenClient creates a new AuthTokenClient.
func NewAuthTokenClient(address, token string) (*AuthTokenClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &AuthTokenClient{
		address: address,
		token:   token,
		client:  &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// ValidateToken checks whether the configured token is valid by calling
// the Vault token lookup-self endpoint.
func (c *AuthTokenClient) ValidateToken() (*TokenInfo, error) {
	url := fmt.Sprintf("%s/v1/auth/token/lookup-self", c.address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token validation request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("token is invalid or expired")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Accessor  string   `json:"accessor"`
			Policies  []string `json:"policies"`
			TTL       int      `json:"ttl"`
			Renewable bool     `json:"renewable"`
			ExpireTime string  `json:"expire_time"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &TokenInfo{
		Accessor:   body.Data.Accessor,
		Policies:   body.Data.Policies,
		TTL:        time.Duration(body.Data.TTL) * time.Second,
		Renewable:  body.Data.Renewable,
		ExpireTime: body.Data.ExpireTime,
	}, nil
}
