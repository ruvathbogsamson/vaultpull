package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// GCPCredentials holds the generated GCP service account credentials.
type GCPCredentials struct {
	Token     string `json:"token"`
	ExpireTime string `json:"expire_time"`
	ServiceAccount string `json:"service_account_email"`
}

// GCPClient interacts with the Vault GCP secrets engine.
type GCPClient struct {
	address string
	token   string
	mount   string
	httpClient *http.Client
}

// NewGCPClient creates a new GCPClient. mount defaults to "gcp" if empty.
func NewGCPClient(address, token, mount string) (*GCPClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	if mount == "" {
		mount = "gcp"
	}
	return &GCPClient{
		address:    strings.TrimRight(address, "/"),
		token:      token,
		mount:      mount,
		httpClient: &http.Client{},
	}, nil
}

// GenerateOAuthToken generates an OAuth2 access token for the given roleset.
func (c *GCPClient) GenerateOAuthToken(roleset string) (*GCPCredentials, error) {
	if roleset == "" {
		return nil, fmt.Errorf("roleset is required")
	}
	url := fmt.Sprintf("%s/v1/%s/token/%s", c.address, c.mount, roleset)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("roleset %q not found", roleset)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var result struct {
		Data GCPCredentials `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result.Data, nil
}
