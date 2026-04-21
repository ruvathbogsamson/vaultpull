// Package vault provides HashiCorp Vault integration for vaultpull.
package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// AppRoleCredentials holds the role_id and secret_id for AppRole auth.
type AppRoleCredentials struct {
	RoleID   string
	SecretID string
}

// AppRoleAuthResult contains the client token returned after a successful login.
type AppRoleAuthResult struct {
	ClientToken   string
	LeaseDuration int
	Renewable     bool
}

// AppRoleClient authenticates against Vault using the AppRole auth method.
type AppRoleClient struct {
	baseURL    string
	mountPath  string
	httpClient *http.Client
}

// NewAppRoleClient creates a new AppRoleClient.
// mountPath defaults to "approle" if empty.
func NewAppRoleClient(baseURL, mountPath string) *AppRoleClient {
	if mountPath == "" {
		mountPath = "approle"
	}
	return &AppRoleClient{
		baseURL:   strings.TrimRight(baseURL, "/"),
		mountPath: mountPath,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Login performs the AppRole login and returns a client token.
func (a *AppRoleClient) Login(ctx context.Context, creds AppRoleCredentials) (*AppRoleAuthResult, error) {
	if creds.RoleID == "" {
		return nil, fmt.Errorf("approle: role_id is required")
	}
	if creds.SecretID == "" {
		return nil, fmt.Errorf("approle: secret_id is required")
	}

	payload := fmt.Sprintf(`{"role_id":%q,"secret_id":%q}`, creds.RoleID, creds.SecretID)
	url := fmt.Sprintf("%s/v1/auth/%s/login", a.baseURL, a.mountPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("approle: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("approle: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("approle: login failed (status %d): %s", resp.StatusCode, body)
	}

	var result struct {
		Auth struct {
			ClientToken   string `json:"client_token"`
			LeaseDuration int    `json:"lease_duration"`
			Renewable     bool   `json:"renewable"`
		} `json:"auth"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("approle: decode response: %w", err)
	}

	return &AppRoleAuthResult{
		ClientToken:   result.Auth.ClientToken,
		LeaseDuration: result.Auth.LeaseDuration,
		Renewable:     result.Auth.Renewable,
	}, nil
}
