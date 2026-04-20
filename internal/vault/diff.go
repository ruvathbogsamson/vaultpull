package vault

import "fmt"

// DiffResult holds the outcome of comparing remote secrets against local env.
type DiffResult struct {
	Added     map[string]string
	Removed   map[string]string
	Changed   map[string]string
	Unchanged map[string]string
}

// Diff compares a set of remote secrets fetched from Vault against the current
// local key/value pairs (e.g. parsed from an existing .env file). It returns a
// DiffResult describing what has been added, removed, changed, or left the same.
func Diff(remote, local map[string]string) DiffResult {
	result := DiffResult{
		Added:     make(map[string]string),
		Removed:   make(map[string]string),
		Changed:   make(map[string]string),
		Unchanged: make(map[string]string),
	}

	for k, rv := range remote {
		lv, exists := local[k]
		if !exists {
			result.Added[k] = rv
		} else if lv != rv {
			result.Changed[k] = rv
		} else {
			result.Unchanged[k] = rv
		}
	}

	for k, lv := range local {
		if _, exists := remote[k]; !exists {
			result.Removed[k] = lv
		}
	}

	return result
}

// HasChanges returns true when there is at least one added, removed, or
// changed key between the remote and local sets.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Summary returns a concise human-readable description of the diff.
func (d DiffResult) Summary() string {
	return fmt.Sprintf("+%d added  ~%d changed  -%d removed  =%d unchanged",
		len(d.Added), len(d.Changed), len(d.Removed), len(d.Unchanged))
}
