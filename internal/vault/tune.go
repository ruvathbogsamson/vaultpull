package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// TuneConfig holds mount tuning parameters for a Vault secrets engine.
type TuneConfig struct {
	DefaultLeaseTTL string `json:"default_lease_ttl"`
	MaxLeaseTTL     string `json:"max_lease_ttl"`
	Description     string `json:"description"`
	ForceNoCache    bool   `json:"force_no_cache"`
}

// TuneClient interacts with Vault mount tuning endpoints.
type TuneClient struct {
	address string
	token   string
	hc      *http.Client
}

// NewTuneClient creates a new TuneClient.
func NewTuneClient(address, token string) (*TuneClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &TuneClient{
		address: strings.TrimRight(address, "/"),
		token:   token,
		hc:      &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// GetTune retrieves tuning configuration for the given mount path.
func (c *TuneClient) GetTune(mount string) (*TuneConfig, error) {
	if mount == "" {
		return nil, fmt.Errorf("mount path is required")
	}
	mount = strings.Trim(mount, "/")
	url := fmt.Sprintf("%s/v1/sys/mounts/%s/tune", c.address, mount)

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
		return nil, fmt.Errorf("mount %q not found", mount)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var cfg TuneConfig
	if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &cfg, nil
}
