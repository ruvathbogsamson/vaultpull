package vault

import (
	"fmt"
	"strings"
)

// TagFilter holds tag-based filtering criteria for secrets.
type TagFilter struct {
	Required map[string]string
}

// NewTagFilter creates a TagFilter from a slice of "key=value" strings.
func NewTagFilter(pairs []string) (*TagFilter, error) {
	required := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("invalid tag filter %q: must be key=value", p)
		}
		required[parts[0]] = parts[1]
	}
	return &TagFilter{Required: required}, nil
}

// Match reports whether the given metadata tags satisfy all required filters.
// An empty TagFilter matches everything.
func (f *TagFilter) Match(tags map[string]string) bool {
	for k, v := range f.Required {
		got, ok := tags[k]
		if !ok || got != v {
			return false
		}
	}
	return true
}

// FilterSecrets returns only those secrets whose custom_metadata satisfies f.
func FilterSecrets(secrets map[string]string, meta map[string]map[string]string, f *TagFilter) map[string]string {
	if f == nil || len(f.Required) == 0 {
		return secrets
	}
	out := make(map[string]string)
	for k, v := range secrets {
		tags := meta[k]
		if f.Match(tags) {
			out[k] = v
		}
	}
	return out
}
