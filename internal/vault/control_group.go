// Package vault provides clients for interacting with HashiCorp Vault.
package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ControlGroupRequest represents a Vault control group authorization request.
type ControlGroupRequest struct {
	WrappingToken string    `json:"wrapping_token"`
	ApprovedBy    []string  `json:"approved_by,omitempty"`
	Approved      bool      `json:"approved"`
	CreatedAt     time.Time `json:"created_at"`
}

// ControlGroupClient interacts with Vault's control group endpoints.
type ControlGroupClient struct {
	address string
	token   string
	httpClient *http.Client
}

// NewControlGroupClient creates a new ControlGroupClient.
// Returns an error if address or token are empty.
func NewControlGroupClient(address, token string) (*ControlGroupClient, error) {
	if strings.TrimSpace(address) == "" {
		return nil, fmt.Errorf("vault address must not be empty")
	}
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("vault token must not be empty")
	}
	return &ControlGroupClient{
		address:    strings.TrimRight(address, "/"),
		token:      token,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Authorize approves a control group request on behalf of the current token.
// wrappingToken is the token associated with the pending request.
func (c *ControlGroupClient) Authorize(wrappingToken string) (*ControlGroupRequest, error) {
	if strings.TrimSpace(wrappingToken) == "" {
		return nil, fmt.Errorf("wrapping token must not be empty")
	}

	body := fmt.Sprintf(`{"token":%q}`, wrappingToken)
	req, err := http.NewRequest(http.MethodPost,
		c.address+"/v1/sys/control-group/authorize",
		strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("building authorize request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("authorize request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("permission denied: token lacks control group authorize capability")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from authorize endpoint", resp.StatusCode)
	}

	var wrapper struct {
		Data ControlGroupRequest `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decoding authorize response: %w", err)
	}
	return &wrapper.Data, nil
}

// Check returns the current status of a control group request.
func (c *ControlGroupClient) Check(wrappingToken string) (*ControlGroupRequest, error) {
	if strings.TrimSpace(wrappingToken) == "" {
		return nil, fmt.Errorf("wrapping token must not be empty")
	}

	body := fmt.Sprintf(`{"token":%q}`, wrappingToken)
	req, err := http.NewRequest(http.MethodPost,
		c.address+"/v1/sys/control-group/request",
		strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("building check request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("check request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("control group request not found for token")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from check endpoint", resp.StatusCode)
	}

	var wrapper struct {
		Data ControlGroupRequest `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decoding check response: %w", err)
	}
	return &wrapper.Data, nil
}
