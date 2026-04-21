package vault

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/api"
)

// EngineType represents a Vault secrets engine type.
type EngineType string

const (
	EngineKV1    EngineType = "kv"
	EngineKV2    EngineType = "kv-v2"
	EngineGeneric EngineType = "generic"
	EngineUnknown EngineType = "unknown"
)

// MountInfo describes a mounted secrets engine.
type MountInfo struct {
	Path        string     `json:"path"`
	Type        EngineType `json:"type"`
	Description string     `json:"description"`
	Version     string     `json:"version"`
}

// EngineClient fetches information about mounted secrets engines.
type EngineClient struct {
	client *api.Client
}

// NewEngineClient creates a new EngineClient.
func NewEngineClient(client *api.Client) *EngineClient {
	return &EngineClient{client: client}
}

// ListMounts returns all mounted secrets engines.
func (e *EngineClient) ListMounts() ([]MountInfo, error) {
	resp, err := e.client.RawRequest(e.client.NewRequest(http.MethodGet, "/v1/sys/mounts"))
	if err != nil {
		return nil, fmt.Errorf("listing mounts: %w", err)
	}
	defer resp.Body.Close()

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding mounts response: %w", err)
	}

	var mounts []MountInfo
	for path, data := range raw {
		var entry struct {
			Type        string `json:"type"`
			Description string `json:"description"`
			Options     struct {
				Version string `json:"version"`
			} `json:"options"`
		}
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}
		if entry.Type == "" {
			continue
		}
		mounts = append(mounts, MountInfo{
			Path:        path,
			Type:        EngineType(entry.Type),
			Description: entry.Description,
			Version:     entry.Options.Version,
		})
	}
	return mounts, nil
}

// GetMount returns info for a specific mount path.
func (e *EngineClient) GetMount(path string) (*MountInfo, error) {
	mounts, err := e.ListMounts()
	if err != nil {
		return nil, err
	}
	for _, m := range mounts {
		if m.Path == path || m.Path == path+"/" {
			return &m, nil
		}
	}
	return nil, fmt.Errorf("mount %q not found", path)
}
