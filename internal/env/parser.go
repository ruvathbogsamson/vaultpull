// Package env provides utilities for reading, writing, and filtering
// environment variable files (.env).
package env

import (
	"bufio"
	"io"
	"strings"
)

// Parse reads key=value pairs from r, ignoring blank lines and lines that
// start with '#'. Inline comments (# after a value) are not stripped so that
// values are preserved verbatim. Surrounding whitespace around the key and
// value is trimmed. Quoted values (single or double) are unquoted.
//
// The returned map contains only the keys present in r; duplicate keys are
// resolved by keeping the last occurrence.
func Parse(r io.Reader) (map[string]string, error) {
	result := make(map[string]string)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			// No '=' found — skip malformed line.
			continue
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		if len(key) == 0 {
			continue
		}

		val = unquote(val)
		result[key] = val
	}

	return result, scanner.Err()
}

// unquote strips a matching pair of leading/trailing single or double quotes
// from s, if present.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
