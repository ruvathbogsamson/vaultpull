package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// QuotaInfo holds rate limit quota details for a Vault path.
type QuotaInfo struct {
	Name          string  `json:"name"`
	Path          string  `json:"path"`
	Type          string  `json:"type"`
	MaxRequests   int     `json:"max_requests"`
	Interval      float64 `json:"interval"`
	Rate          float64 `json:"rate"`
	BlockInterval float64 `json:"block_interval"`
}

// QuotaClient interacts with Vault's quota API.
type QuotaClient struct {
	address string
	token   string
	hc      *http.Client
}

// NewQuotaClient creates a new QuotaClient.
func NewQuotaClient(address, token string) (*QuotaClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &QuotaClient{
		address: address,
		token:   token,
		hc:      &http.Client{},
	}, nil
}

// GetQuota fetches a named rate-limit quota from Vault.
func (q *QuotaClient) GetQuota(name string) (*QuotaInfo, error) {
	if name == "" {
		return nil, fmt.Errorf("quota name is required")
	}
	url := fmt.Sprintf("%s/v1/sys/quotas/rate-limit/%s", q.address, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", q.token)

	resp, err := q.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("quota %q not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var envelope struct {
		Data QuotaInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &envelope.Data, nil
}
