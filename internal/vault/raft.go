package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RaftPeer represents a single member of a Vault Raft cluster.
type RaftPeer struct {
	NodeID   string `json:"node_id"`
	Address  string `json:"address"`
	Leader   bool   `json:"leader"`
	Protocol string `json:"protocol_version"`
	Voter    bool   `json:"voter"`
}

// RaftConfig holds the Raft cluster configuration returned by Vault.
type RaftConfig struct {
	Index   uint64      `json:"index"`
	Servers []RaftPeer  `json:"servers"`
}

// RaftClient interacts with the Vault Raft operator endpoints.
type RaftClient struct {
	address string
	token   string
	hc      *http.Client
}

// NewRaftClient creates a new RaftClient. Returns an error if address or token
// is empty.
func NewRaftClient(address, token string) (*RaftClient, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required")
	}
	return &RaftClient{
		address: address,
		token:   token,
		hc:      &http.Client{},
	}, nil
}

// GetRaftConfig fetches the current Raft cluster configuration from Vault.
func (c *RaftClient) GetRaftConfig() (*RaftConfig, error) {
	url := fmt.Sprintf("%s/v1/sys/storage/raft/configuration", c.address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("raft configuration not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wrapper struct {
		Data RaftConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &wrapper.Data, nil
}
