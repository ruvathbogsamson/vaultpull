package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// StepDownClient interacts with the Vault operator step-down endpoint.
type StepDownClient struct {
	address string
	token   string
	httpClient *http.Client
}

// StepDownStatus holds the result of a step-down operation.
type StepDownStatus struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// NewStepDownClient creates a new StepDownClient.
func NewStepDownClient(address, token string) (*StepDownClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &StepDownClient{
		address:    address,
		token:      token,
		httpClient: &http.Client{},
	}, nil
}

// StepDown forces the active Vault node to step down as leader.
func (c *StepDownClient) StepDown() (*StepDownStatus, error) {
	url := fmt.Sprintf("%s/v1/sys/step-down", c.address)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating step-down request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("step-down request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		return &StepDownStatus{Success: true, Message: "step-down successful"}, nil
	}

	var apiErr map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil {
		if errs, ok := apiErr["errors"].([]interface{}); ok && len(errs) > 0 {
			return nil, fmt.Errorf("step-down failed: %v", errs[0])
		}
	}
	return nil, fmt.Errorf("step-down failed with status %d", resp.StatusCode)
}
