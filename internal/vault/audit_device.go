package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AuditDevice represents a Vault audit device configuration.
type AuditDevice struct {
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Options     map[string]string `json:"options"`
	Path        string            `json:"path"`
}

// AuditDeviceClient interacts with Vault audit device endpoints.
type AuditDeviceClient struct {
	address string
	token   string
	client  *http.Client
}

// NewAuditDeviceClient creates a new AuditDeviceClient.
func NewAuditDeviceClient(address, token string) (*AuditDeviceClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &AuditDeviceClient{
		address: address,
		token:   token,
		client:  &http.Client{},
	}, nil
}

// ListAuditDevices returns all enabled audit devices.
func (c *AuditDeviceClient) ListAuditDevices() (map[string]*AuditDevice, error) {
	req, err := http.NewRequest(http.MethodGet, c.address+"/v1/sys/audit", nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("audit endpoint not found")
	}
	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("permission denied")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result map[string]*AuditDevice
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return result, nil
}
