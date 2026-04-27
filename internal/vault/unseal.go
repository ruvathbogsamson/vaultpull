package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// UnsealClient interacts with the Vault unseal API.
type UnsealClient struct {
	address string
	token   string
	client  *http.Client
}

// UnsealRequest represents the payload for submitting an unseal key.
type UnsealRequest struct {
	Key   string `json:"key"`
	Reset bool   `json:"reset,omitempty"`
}

// UnsealResponse contains the current seal status after submitting a key.
type UnsealResponse struct {
	Sealed   bool `json:"sealed"`
	T        int  `json:"t"`
	N        int  `json:"n"`
	Progress int  `json:"progress"`
}

// NewUnsealClient creates a new UnsealClient.
func NewUnsealClient(address, token string) (*UnsealClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &UnsealClient{
		address: address,
		token:   token,
		client:  &http.Client{},
	}, nil
}

// SubmitKey submits a single unseal key shard to Vault.
func (u *UnsealClient) SubmitKey(key string) (*UnsealResponse, error) {
	payload, err := json.Marshal(UnsealRequest{Key: key})
	if err != nil {
		return nil, fmt.Errorf("marshal unseal request: %w", err)
	}

	req, err := newJSONRequest(http.MethodPut, u.address+"/v1/sys/unseal", u.token, payload)
	if err != nil {
		return nil, fmt.Errorf("build unseal request: %w", err)
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unseal request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result UnsealResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode unseal response: %w", err)
	}
	return &result, nil
}

// Reset cancels an in-progress unseal attempt.
func (u *UnsealClient) Reset() (*UnsealResponse, error) {
	payload, err := json.Marshal(UnsealRequest{Reset: true})
	if err != nil {
		return nil, fmt.Errorf("marshal reset request: %w", err)
	}

	req, err := newJSONRequest(http.MethodPut, u.address+"/v1/sys/unseal", u.token, payload)
	if err != nil {
		return nil, fmt.Errorf("build reset request: %w", err)
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("reset request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result UnsealResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode reset response: %w", err)
	}
	return &result, nil
}
