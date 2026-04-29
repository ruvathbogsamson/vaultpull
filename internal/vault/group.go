package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Group represents a Vault identity group.
type Group struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Metadata map[string]string `json:"metadata"`
	MemberEntityIDs []string   `json:"member_entity_ids"`
	Policies []string          `json:"policies"`
}

// GroupClient interacts with the Vault identity group API.
type GroupClient struct {
	address string
	token   string
	client  *http.Client
}

// NewGroupClient creates a new GroupClient.
func NewGroupClient(address, token string) (*GroupClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &GroupClient{
		address: address,
		token:   token,
		client:  &http.Client{},
	}, nil
}

// GetGroup retrieves a Vault identity group by name.
func (c *GroupClient) GetGroup(name string) (*Group, error) {
	if name == "" {
		return nil, fmt.Errorf("group name is required")
	}
	url := fmt.Sprintf("%s/v1/identity/group/name/%s", c.address, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("group %q not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Data Group `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result.Data, nil
}
