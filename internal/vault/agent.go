package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// AgentConfig holds the result of reading a Vault Agent auto-auth configuration.
type AgentConfig struct {
	Token     string            `json:"token"`
	Mount     string            `json:"mount"`
	Metadata  map[string]string `json:"metadata"`
}

// AgentClient reads token and metadata from a running Vault Agent.
type AgentClient struct {
	address string
	hc      *http.Client
}

// NewAgentClient creates an AgentClient targeting the given Vault Agent address.
// Returns an error if address is empty.
func NewAgentClient(address string) (*AgentClient, error) {
	address = strings.TrimRight(address, "/")
	if address == "" {
		return nil, fmt.Errorf("agent address must not be empty")
	}
	return &AgentClient{
		address: address,
		hc:      &http.Client{},
	}, nil
}

// FetchConfig retrieves the current agent token and metadata from the agent
// sink endpoint at /v1/agent/self.
func (c *AgentClient) FetchConfig() (*AgentConfig, error) {
	url := c.address + "/v1/agent/self"
	resp, err := c.hc.Get(url) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("agent request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("agent endpoint not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("agent returned status %d", resp.StatusCode)
	}

	var wrapper struct {
		Data AgentConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("failed to decode agent response: %w", err)
	}
	return &wrapper.Data, nil
}

// IsHealthy checks whether the Vault Agent is reachable and responding at its
// health endpoint (/v1/agent/health). Returns nil if the agent is healthy.
func (c *AgentClient) IsHealthy() error {
	url := c.address + "/v1/agent/health"
	resp, err := c.hc.Get(url) //nolint:noctx
	if err != nil {
		return fmt.Errorf("agent health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("agent health check returned status %d", resp.StatusCode)
	}
	return nil
}
