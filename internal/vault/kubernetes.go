package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// KubernetesClient handles Vault Kubernetes auth method login.
type KubernetesClient struct {
	address string
	mount   string
	client  *http.Client
}

// KubernetesLoginRequest holds the payload for a Kubernetes login.
type KubernetesLoginRequest struct {
	Role string `json:"role"`
	JWT  string `json:"jwt"`
}

// KubernetesLoginResponse holds the Vault token returned after login.
type KubernetesLoginResponse struct {
	ClientToken string
	Accessor    string
	Policies    []string
}

// NewKubernetesClient creates a new KubernetesClient.
// mount defaults to "kubernetes" if empty.
func NewKubernetesClient(address, mount string) (*KubernetesClient, error) {
	if strings.TrimSpace(address) == "" {
		return nil, fmt.Errorf("kubernetes: vault address is required")
	}
	if mount == "" {
		mount = "kubernetes"
	}
	return &KubernetesClient{
		address: strings.TrimRight(address, "/"),
		mount:   mount,
		client:  &http.Client{},
	}, nil
}

// Login authenticates using the Kubernetes auth method and returns a login response.
func (k *KubernetesClient) Login(role, jwt string) (*KubernetesLoginResponse, error) {
	if strings.TrimSpace(role) == "" {
		return nil, fmt.Errorf("kubernetes: role is required")
	}
	if strings.TrimSpace(jwt) == "" {
		return nil, fmt.Errorf("kubernetes: jwt is required")
	}

	payload, err := json.Marshal(KubernetesLoginRequest{Role: role, JWT: jwt})
	if err != nil {
		return nil, fmt.Errorf("kubernetes: failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/auth/%s/login", k.address, k.mount)
	resp, err := k.client.Post(url, "application/json", strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("kubernetes: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("kubernetes: authentication failed (status %d)", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kubernetes: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string   `json:"client_token"`
			Accessor    string   `json:"accessor"`
			Policies    []string `json:"policies"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("kubernetes: failed to decode response: %w", err)
	}

	return &KubernetesLoginResponse{
		ClientToken: result.Auth.ClientToken,
		Accessor:    result.Auth.Accessor,
		Policies:    result.Auth.Policies,
	}, nil
}
