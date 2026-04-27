package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	hashivault "github.com/hashicorp/vault/api"
)

// MountInfo holds configuration details for a single secrets engine mount.
type MountInfo struct {
	Path        string
	Type        string
	Description string
	Options     map[string]string
}

// MountClient interacts with Vault's sys/mounts endpoint.
type MountClient struct {
	client *hashivault.Client
}

// NewMountClient creates a MountClient. Returns an error if address or token
// are empty.
func NewMountClient(address, token string) (*MountClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	cfg := hashivault.DefaultConfig()
	cfg.Address = address
	c, err := hashivault.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}
	c.SetToken(token)
	return &MountClient{client: c}, nil
}

// ListMountPaths returns all mount paths visible to the token.
func (m *MountClient) ListMountPaths() ([]MountInfo, error) {
	resp, err := m.client.RawRequest(m.client.NewRequest(http.MethodGet, "/v1/sys/mounts"))
	if err != nil {
		return nil, fmt.Errorf("listing mounts: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	var raw map[string]struct {
		Type        string            `json:"type"`
		Description string            `json:"description"`
		Options     map[string]string `json:"options"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding mounts response: %w", err)
	}
	var mounts []MountInfo
	for path, info := range raw {
		mounts = append(mounts, MountInfo{
			Path:        strings.TrimSuffix(path, "/"),
			Type:        info.Type,
			Description: info.Description,
			Options:     info.Options,
		})
	}
	return mounts, nil
}
