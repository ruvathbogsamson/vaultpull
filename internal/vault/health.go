package vault

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// HealthStatus represents the result of a Vault health check.
type HealthStatus struct {
	Initialized bool
	Sealed      bool
	Standby     bool
	Version     string
	Reachable   bool
}

// CheckHealth performs a health check against the Vault server.
// It uses the /v1/sys/health endpoint which returns status without auth.
func (c *Client) CheckHealth(ctx context.Context) (*HealthStatus, error) {
	url := fmt.Sprintf("%s/v1/sys/health", c.address)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building health request: %w", err)
	}

	httpClient := &http.Client{Timeout: 5 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return &HealthStatus{Reachable: false}, fmt.Errorf("vault unreachable: %w", err)
	}
	defer resp.Body.Close()

	status := &HealthStatus{Reachable: true}

	switch resp.StatusCode {
	case http.StatusOK:
		status.Initialized = true
		status.Sealed = false
		status.Standby = false
	case http.StatusTooManyRequests:
		status.Initialized = true
		status.Sealed = false
		status.Standby = true
	case http.StatusNotImplemented:
		status.Initialized = false
	case 503:
		status.Initialized = true
		status.Sealed = true
	default:
		return status, fmt.Errorf("unexpected health status code: %d", resp.StatusCode)
	}

	return status, nil
}
