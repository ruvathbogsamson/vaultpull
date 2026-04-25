package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SealStatus represents the current seal state of a Vault instance.
type SealStatus struct {
	Sealed      bool   `json:"sealed"`
	Initialized bool   `json:"initialized"`
	T           int    `json:"t"`
	N           int    `json:"n"`
	Progress    int    `json:"progress"`
	Version     string `json:"version"`
	ClusterName string `json:"cluster_name"`
	ClusterID   string `json:"cluster_id"`
}

// SealClient interacts with Vault seal/unseal endpoints.
type SealClient struct {
	address string
	token   string
	httpClient *http.Client
}

// NewSealClient creates a new SealClient.
func NewSealClient(address, token string) (*SealClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &SealClient{
		address:    address,
		token:      token,
		httpClient: &http.Client{},
	}, nil
}

// GetSealStatus returns the current seal status of the Vault.
func (c *SealClient) GetSealStatus() (*SealStatus, error) {
	url := fmt.Sprintf("%s/v1/sys/seal-status", c.address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var status SealStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &status, nil
}

// Seal seals the Vault.
func (c *SealClient) Seal() error {
	url := fmt.Sprintf("%s/v1/sys/seal", c.address)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}
