package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AWSCredentials holds the temporary AWS credentials returned by Vault.
type AWSCredentials struct {
	AccessKey     string `json:"access_key"`
	SecretKey     string `json:"secret_key"`
	SecurityToken string `json:"security_token"`
	LeaseDuration int    `json:"lease_duration"`
}

// AWSClient fetches dynamic AWS credentials from Vault's AWS secrets engine.
type AWSClient struct {
	address string
	token   string
	mount   string
	client  *http.Client
}

// NewAWSClient creates a new AWSClient. mount defaults to "aws" if empty.
func NewAWSClient(address, token, mount string) (*AWSClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	if mount == "" {
		mount = "aws"
	}
	return &AWSClient{
		address: address,
		token:   token,
		mount:   mount,
		client:  &http.Client{},
	}, nil
}

// GenerateCredentials requests dynamic AWS credentials for the given role.
func (c *AWSClient) GenerateCredentials(role string) (*AWSCredentials, error) {
	if role == "" {
		return nil, fmt.Errorf("role is required")
	}
	url := fmt.Sprintf("%s/v1/%s/creds/%s", c.address, c.mount, role)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting credentials: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("role %q not found", role)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var envelope struct {
		Data          AWSCredentials `json:"data"`
		LeaseDuration int            `json:"lease_duration"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	envelope.Data.LeaseDuration = envelope.LeaseDuration
	return &envelope.Data, nil
}
