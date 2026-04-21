package vault

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TransitClient provides encrypt/decrypt operations via Vault's Transit engine.
type TransitClient struct {
	address string
	token   string
	mount   string
	httpClient *http.Client
}

// NewTransitClient creates a TransitClient for the given mount (defaults to "transit").
func NewTransitClient(address, token, mount string) (*TransitClient, error) {
	if address == "" {
		return nil, fmt.Errorf("transit: vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("transit: vault token is required")
	}
	if mount == "" {
		mount = "transit"
	}
	return &TransitClient{
		address:    strings.TrimRight(address, "/"),
		token:      token,
		mount:      mount,
		httpClient: &http.Client{},
	}, nil
}

// Encrypt encrypts plaintext using the named key and returns base64 ciphertext.
func (c *TransitClient) Encrypt(keyName, plaintext string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(plaintext))
	body, _ := json.Marshal(map[string]string{"plaintext": encoded})
	url := fmt.Sprintf("%s/v1/%s/encrypt/%s", c.address, c.mount, keyName)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Vault-Token", c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("transit: encrypt request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("transit: encrypt returned status %d", resp.StatusCode)
	}
	var result struct {
		Data struct {
			Ciphertext string `json:"ciphertext"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("transit: failed to decode encrypt response: %w", err)
	}
	return result.Data.Ciphertext, nil
}

// Decrypt decrypts ciphertext using the named key and returns plaintext.
func (c *TransitClient) Decrypt(keyName, ciphertext string) (string, error) {
	body, _ := json.Marshal(map[string]string{"ciphertext": ciphertext})
	url := fmt.Sprintf("%s/v1/%s/decrypt/%s", c.address, c.mount, keyName)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Vault-Token", c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("transit: decrypt request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("transit: decrypt returned status %d: %s", resp.StatusCode, string(b))
	}
	var result struct {
		Data struct {
			Plaintext string `json:"plaintext"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("transit: failed to decode decrypt response: %w", err)
	}
	decoded, err := base64.StdEncoding.DecodeString(result.Data.Plaintext)
	if err != nil {
		return "", fmt.Errorf("transit: failed to base64-decode plaintext: %w", err)
	}
	return string(decoded), nil
}
