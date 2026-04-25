package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RekeyStatus holds the current status of a rekey operation.
type RekeyStatus struct {
	Started        bool     `json:"started"`
	T              int      `json:"t"`
	N              int      `json:"n"`
	Progress       int      `json:"progress"`
	Required       int      `json:"required"`
	PGPFingerprints []string `json:"pgp_fingerprints"`
	Backup         bool     `json:"backup"`
	Nonce          string   `json:"nonce"`
}

// RekeyClient interacts with Vault's rekey API.
type RekeyClient struct {
	address string
	token   string
	hc      *http.Client
}

// NewRekeyClient creates a new RekeyClient.
func NewRekeyClient(address, token string) (*RekeyClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &RekeyClient{
		address: address,
		token:   token,
		hc:      &http.Client{},
	}, nil
}

// GetRekeyStatus returns the current rekey status from Vault.
func (c *RekeyClient) GetRekeyStatus() (*RekeyStatus, error) {
	url := fmt.Sprintf("%s/v1/sys/rekey/init", c.address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("rekey endpoint not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var status RekeyStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &status, nil
}
