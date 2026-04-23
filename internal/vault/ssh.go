package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SSHCredential holds the signed key or OTP returned by Vault.
type SSHCredential struct {
	SignedKey string `json:"signed_key"`
	OTP       string `json:"key"`
	KeyType   string `json:"key_type"`
}

// SSHClient interacts with the Vault SSH secrets engine.
type SSHClient struct {
	address string
	token   string
	mount   string
	client  *http.Client
}

// NewSSHClient creates a new SSHClient.
// address and token are required; mount defaults to "ssh".
func NewSSHClient(address, token, mount string) (*SSHClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	if mount == "" {
		mount = "ssh"
	}
	return &SSHClient{
		address: address,
		token:   token,
		mount:   mount,
		client:  &http.Client{},
	}, nil
}

// SignKey requests Vault to sign the provided public key for the given role.
func (c *SSHClient) SignKey(role, publicKey string) (*SSHCredential, error) {
	if role == "" {
		return nil, fmt.Errorf("role is required")
	}
	if publicKey == "" {
		return nil, fmt.Errorf("public key is required")
	}

	url := fmt.Sprintf("%s/v1/%s/sign/%s", c.address, c.mount, role)
	body, _ := json.Marshal(map[string]string{"public_key": publicKey})

	req, err := http.NewRequest(http.MethodPost, url, jsonReader(body))
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("role %q not found", role)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var wrapper struct {
		Data SSHCredential `json:"data"`
	}
	if err := json.Unmarshal(raw, &wrapper); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &wrapper.Data, nil
}
