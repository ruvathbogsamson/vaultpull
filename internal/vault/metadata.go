package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SecretMetadata holds KV v2 metadata for a secret path.
type SecretMetadata struct {
	CreatedTime    time.Time         `json:"created_time"`
	UpdatedTime    time.Time         `json:"updated_time"`
	CurrentVersion int               `json:"current_version"`
	OldestVersion  int               `json:"oldest_version"`
	CustomMetadata map[string]string `json:"custom_metadata"`
}

// MetadataClient fetches KV v2 secret metadata from Vault.
type MetadataClient struct {
	httpClient *http.Client
	address    string
	token      string
	mount      string
}

// NewMetadataClient creates a MetadataClient for the given mount.
func NewMetadataClient(address, token, mount string) *MetadataClient {
	return &MetadataClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		address:    address,
		token:      token,
		mount:      mount,
	}
}

// FetchMetadata retrieves metadata for the secret at path.
func (m *MetadataClient) FetchMetadata(path string) (*SecretMetadata, error) {
	url := fmt.Sprintf("%s/v1/%s/metadata/%s", m.address, m.mount, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("metadata: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", m.token)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("metadata: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("metadata: path %q not found", path)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("metadata: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			CreatedTime    time.Time         `json:"created_time"`
			UpdatedTime    time.Time         `json:"updated_time"`
			CurrentVersion int               `json:"current_version"`
			OldestVersion  int               `json:"oldest_version"`
			CustomMetadata map[string]string `json:"custom_metadata"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("metadata: decode response: %w", err)
	}

	return &SecretMetadata{
		CreatedTime:    body.Data.CreatedTime,
		UpdatedTime:    body.Data.UpdatedTime,
		CurrentVersion: body.Data.CurrentVersion,
		OldestVersion:  body.Data.OldestVersion,
		CustomMetadata: body.Data.CustomMetadata,
	}, nil
}
