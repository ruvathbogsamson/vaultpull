package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// TOTPClient interacts with the Vault TOTP secrets engine to generate
// and validate time-based one-time passwords.
type TOTPClient struct {
	address string
	token   string
	mount   string
	httpClient *http.Client
}

// TOTPCode holds a generated TOTP code returned by Vault.
type TOTPCode struct {
	Code      string    `json:"code"`
	GeneratedAt time.Time `json:"-"`
}

// TOTPValidateResult holds the result of a TOTP validation request.
type TOTPValidateResult struct {
	Valid bool `json:"valid"`
}

totpDefaultMount = "totp"

// NewTOTPClient creates a new TOTPClient. address and token are required;
// mount defaults to "totp" if empty.
func NewTOTPClient(address, token, mount string) (*TOTPClient, error) {
	if strings.TrimSpace(address) == "" {
		return nil, fmt.Errorf("totp: vault address is required")
	}
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("totp: vault token is required")
	}
	if mount == "" {
		mount = totpDefaultMount
	}
	return &TOTPClient{
		address:    strings.TrimRight(address, "/"),
		token:      token,
		mount:      mount,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// GenerateCode requests a TOTP code for the named key from Vault.
func (c *TOTPClient) GenerateCode(keyName string) (*TOTPCode, error) {
	if keyName == "" {
		return nil, fmt.Errorf("totp: key name is required")
	}
	url := fmt.Sprintf("%s/v1/%s/code/%s", c.address, c.mount, keyName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("totp: building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("totp: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("totp: key %q not found", keyName)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("totp: unexpected status %d", resp.StatusCode)
	}

	var payload struct {
		Data TOTPCode `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("totp: decoding response: %w", err)
	}
	payload.Data.GeneratedAt = time.Now().UTC()
	return &payload.Data, nil
}

// ValidateCode validates a TOTP code for the named key against Vault.
func (c *TOTPClient) ValidateCode(keyName, code string) (*TOTPValidateResult, error) {
	if keyName == "" {
		return nil, fmt.Errorf("totp: key name is required")
	}
	if code == "" {
		return nil, fmt.Errorf("totp: code is required")
	}

	body, err := json.Marshal(map[string]string{"code": code})
	if err != nil {
		return nil, fmt.Errorf("totp: marshalling request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/%s/code/%s", c.address, c.mount, keyName)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("totp: building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("totp: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("totp: key %q not found", keyName)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("totp: unexpected status %d", resp.StatusCode)
	}

	var payload struct {
		Data TOTPValidateResult `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("totp: decoding response: %w", err)
	}
	return &payload.Data, nil
}
