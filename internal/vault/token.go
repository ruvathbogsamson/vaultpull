package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TokenInfo holds metadata about the current Vault token.
type TokenInfo struct {
	Accessor   string            `json:"accessor"`
	Policies   []string          `json:"policies"`
	TTL        time.Duration     `json:"-"`
	TTLSeconds int               `json:"ttl"`
	Renewable  bool              `json:"renewable"`
	Meta       map[string]string `json:"meta"`
	DisplayName string           `json:"display_name"`
}

// TokenClient provides operations for inspecting the current Vault token.
type TokenClient struct {
	httpClient *http.Client
	address    string
	token      string
}

// NewTokenClient creates a TokenClient using the provided Vault client fields.
func NewTokenClient(address, token string, httpClient *http.Client) *TokenClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &TokenClient{
		httpClient: httpClient,
		address:    address,
		token:      token,
	}
}

// LookupSelf retrieves metadata about the token currently in use.
func (tc *TokenClient) LookupSelf() (*TokenInfo, error) {
	url := fmt.Sprintf("%s/v1/auth/token/lookup-self", tc.address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("token: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", tc.token)

	resp, err := tc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("token: permission denied (status %d)", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data TokenInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("token: decode response: %w", err)
	}

	body.Data.TTL = time.Duration(body.Data.TTLSeconds) * time.Second
	return &body.Data, nil
}
