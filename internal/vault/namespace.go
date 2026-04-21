package vault

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

// NamespaceClient provides operations scoped to a Vault namespace.
type NamespaceClient struct {
	client *api.Client
	namespace string
}

// NewNamespaceClient creates a NamespaceClient scoped to the given namespace.
// namespace may be a slash-separated path (e.g. "team/project").
func NewNamespaceClient(addr, token, namespace string) (*NamespaceClient, error) {
	if addr == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	cfg := api.DefaultConfig()
	cfg.Address = addr

	c, err := api.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}
	c.SetToken(token)

	// Vault namespace header must not have a leading slash.
	ns := strings.TrimPrefix(namespace, "/")
	c.SetNamespace(ns)

	return &NamespaceClient{client: c, namespace: ns}, nil
}

// Namespace returns the active namespace.
func (n *NamespaceClient) Namespace() string {
	return n.namespace
}

// ListNamespaces returns child namespaces under the current namespace.
func (n *NamespaceClient) ListNamespaces() ([]string, error) {
	secret, err := n.client.Logical().List("sys/namespaces")
	if err != nil {
		return nil, fmt.Errorf("listing namespaces: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return []string{}, nil
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return []string{}, nil
	}

	result := make([]string, 0, len(keys))
	for _, k := range keys {
		if s, ok := k.(string); ok {
			result = append(result, strings.TrimSuffix(s, "/"))
		}
	}
	return result, nil
}
