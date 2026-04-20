package vault

import (
	"encoding/json"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Snapshot represents a point-in-time capture of secrets at a given path.
type Snapshot struct {
	Path      string            `json:"path"`
	Data      map[string]string `json:"data"`
	CapturedAt time.Time        `json:"captured_at"`
}

// SnapshotClient captures and restores secret snapshots from Vault.
type SnapshotClient struct {
	client *vaultapi.Client
}

// NewSnapshotClient returns a new SnapshotClient using the given Vault API client.
func NewSnapshotClient(c *vaultapi.Client) *SnapshotClient {
	return &SnapshotClient{client: c}
}

// Capture reads secrets at path and returns a Snapshot.
func (s *SnapshotClient) Capture(path string) (*Snapshot, error) {
	secret, err := s.client.Logical().Read(dataPathFromPath(path))
	if err != nil {
		return nil, fmt.Errorf("snapshot capture: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("snapshot capture: no data at path %q", path)
	}

	data, err := flattenSecretData(secret.Data)
	if err != nil {
		return nil, fmt.Errorf("snapshot capture: %w", err)
	}

	return &Snapshot{
		Path:       path,
		Data:       data,
		CapturedAt: time.Now().UTC(),
	}, nil
}

// Marshal serialises a Snapshot to JSON bytes.
func (s *Snapshot) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalSnapshot deserialises JSON bytes into a Snapshot.
func UnmarshalSnapshot(b []byte) (*Snapshot, error) {
	var snap Snapshot
	if err := json.Unmarshal(b, &snap); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot: %w", err)
	}
	return &snap, nil
}

// flattenSecretData converts the raw secret data map (which may contain a
// nested "data" key for KV v2) into a flat map[string]string.
func flattenSecretData(raw map[string]interface{}) (map[string]string, error) {
	src := raw
	if nested, ok := raw["data"]; ok {
		if m, ok := nested.(map[string]interface{}); ok {
			src = m
		}
	}
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out, nil
}
