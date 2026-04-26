package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// MFAClient interacts with Vault's MFA validation endpoints.
type MFAClient struct {
	address string
	token   string
	client  *http.Client
}

// MFAValidateRequest holds the payload for an MFA validation request.
type MFAValidateRequest struct {
	MFARequestID string            `json:"mfa_request_id"`
	MFAPayload   map[string]string `json:"mfa_payload"`
}

// MFAValidateResponse holds the response from a successful MFA validation.
type MFAValidateResponse struct {
	Token    string `json:"client_token"`
	Policies []string
}

// NewMFAClient creates a new MFAClient. Returns an error if address or token is empty.
func NewMFAClient(address, token string) (*MFAClient, error) {
	if address == "" {
		return nil, fmt.Errorf("mfa: vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("mfa: vault token is required")
	}
	return &MFAClient{
		address: address,
		token:   token,
		client:  &http.Client{},
	}, nil
}

// Validate submits an MFA validation request and returns the resulting client token.
func (m *MFAClient) Validate(req MFAValidateRequest) (*MFAValidateResponse, error) {
	if req.MFARequestID == "" {
		return nil, fmt.Errorf("mfa: mfa_request_id is required")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("mfa: failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/sys/mfa/validate", m.address)
	httpReq, err := newJSONRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("mfa: failed to build request: %w", err)
	}
	httpReq.Header.Set("X-Vault-Token", m.token)

	resp, err := m.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("mfa: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("mfa: validation failed: forbidden")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mfa: unexpected status %d", resp.StatusCode)
	}

	var wrapper struct {
		Auth struct {
			ClientToken string   `json:"client_token"`
			Policies    []string `json:"policies"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("mfa: failed to decode response: %w", err)
	}

	return &MFAValidateResponse{
		Token:    wrapper.Auth.ClientToken,
		Policies: wrapper.Auth.Policies,
	}, nil
}
