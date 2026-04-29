package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Entity represents a Vault identity entity.
type Entity struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Policies []string          `json:"policies"`
	Metadata map[string]string `json:"metadata"`
	Disabled bool              `json:"disabled"`
}

// EntityClient interacts with the Vault identity entity API.
type EntityClient struct {
	address string
	token   string
	client  *http.Client
}

// NewEntityClient creates a new EntityClient.
func NewEntityClient(address, token string) (*EntityClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &EntityClient{
		address: address,
		token:   token,
		client:  &http.Client{},
	}, nil
}

// GetEntity retrieves an entity by name from Vault identity store.
func (c *EntityClient) GetEntity(name string) (*Entity, error) {
	if name == "" {
		return nil, fmt.Errorf("entity name is required")
	}
	url := fmt.Sprintf("%s/v1/identity/entity/name/%s", c.address, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting entity: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("entity %q not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wrapper struct {
		Data Entity `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &wrapper.Data, nil
}
