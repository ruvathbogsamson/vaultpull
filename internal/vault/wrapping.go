package vault

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/api"
)

// WrappingClient handles Vault response wrapping (cubbyhole tokens).
type WrappingClient struct {
	client *api.Client
}

// WrappedSecret holds the wrapping token and associated metadata.
type WrappedSecret struct {
	Token    string
	Accessor string
	TTL      int
	Creation string
}

// UnwrappedData holds the key/value pairs from an unwrapped token.
type UnwrappedData map[string]string

// NewWrappingClient creates a new WrappingClient.
func NewWrappingClient(address, token string) (*WrappingClient, error) {
	if strings.TrimSpace(address) == "" {
		return nil, errors.New("vault address is required")
	}
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("vault token is required")
	}
	cfg := api.DefaultConfig()
	cfg.Address = address
	c, err := api.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}
	c.SetToken(token)
	return &WrappingClient{client: c}, nil
}

// Unwrap takes a wrapping token and returns the unwrapped key/value data.
func (w *WrappingClient) Unwrap(wrappingToken string) (UnwrappedData, error) {
	if strings.TrimSpace(wrappingToken) == "" {
		return nil, errors.New("wrapping token is required")
	}
	secret, err := w.client.Logical().Unwrap(wrappingToken)
	if err != nil {
		return nil, fmt.Errorf("unwrapping token: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("no data returned from unwrap")
	}
	result := make(UnwrappedData)
	for k, v := range secret.Data {
		if s, ok := v.(string); ok {
			result[k] = s
		}
	}
	return result, nil
}

// LookupWrappingToken returns metadata about a wrapping token without consuming it.
func (w *WrappingClient) LookupWrappingToken(wrappingToken string) (*WrappedSecret, error) {
	if strings.TrimSpace(wrappingToken) == "" {
		return nil, errors.New("wrapping token is required")
	}
	req := w.client.NewRequest(http.MethodPost, "/v1/sys/wrapping/lookup")
	if err := req.SetJSONBody(map[string]string{"token": wrappingToken}); err != nil {
		return nil, fmt.Errorf("building lookup request: %w", err)
	}
	resp, err := w.client.RawRequest(req)
	if err != nil {
		return nil, fmt.Errorf("lookup request failed: %w", err)
	}
	defer resp.Body.Close()
	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing lookup response: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("empty lookup response")
	}
	ws := &WrappedSecret{}
	if v, ok := secret.Data["token"].(string); ok {
		ws.Token = v
	}
	if v, ok := secret.Data["accessor"].(string); ok {
		ws.Accessor = v
	}
	if v, ok := secret.Data["creation_time"].(string); ok {
		ws.Creation = v
	}
	return ws, nil
}
