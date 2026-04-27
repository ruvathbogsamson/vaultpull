package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// InitRequest holds parameters for initializing a Vault cluster.
type InitRequest struct {
	SecretShares    int `json:"secret_shares"`
	SecretThreshold int `json:"secret_threshold"`
}

// InitResponse holds the unseal keys and root token returned after init.
type InitResponse struct {
	Keys       []string `json:"keys"`
	RootToken  string   `json:"root_token"`
}

// InitClient interacts with the Vault sys/init endpoint.
type InitClient struct {
	address    string
	httpClient *http.Client
}

// NewInitClient creates a new InitClient. No token is required for init.
func NewInitClient(address string) (*InitClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address must not be empty")
	}
	return &InitClient{
		address:    address,
		httpClient: &http.Client{},
	}, nil
}

// Initialize sends an init request to Vault and returns unseal keys and root token.
func (c *InitClient) Initialize(req InitRequest) (*InitResponse, error) {
	if req.SecretShares <= 0 {
		return nil, fmt.Errorf("secret_shares must be greater than zero")
	}
	if req.SecretThreshold <= 0 || req.SecretThreshold > req.SecretShares {
		return nil, fmt.Errorf("secret_threshold must be between 1 and secret_shares")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal init request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/sys/init", c.address)
	httpReq, err := http.NewRequest(http.MethodPost, url, bytesReader(body))
	if err != nil {
		return nil, fmt.Errorf("build init request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("init request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("init returned status %d", resp.StatusCode)
	}

	var result InitResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode init response: %w", err)
	}
	return &result, nil
}
