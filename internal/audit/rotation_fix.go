package audit

import (
	"path/filepath"
	"sort"
)

// sortedBackups returns a sorted list of backup files matching the given
// glob pattern. Files are sorted lexicographically, which works correctly
// for timestamp-based filenames (e.g. audit-2024-01-15T10:30:00.log).
func sortedBackups(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)
	return matches, nil
}
