package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer writes secrets to a .env file.
type Writer struct {
	filePath string
}

// NewWriter creates a new Writer for the given file path.
func NewWriter(filePath string) *Writer {
	return &Writer{filePath: filePath}
}

// Write serializes the given key-value map to the .env file.
// Keys are written in sorted order for deterministic output.
// Values containing spaces or special characters are quoted.
func (w *Writer) Write(secrets map[string]string) error {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, quoteValue(secrets[k])))
	}

	return os.WriteFile(w.filePath, []byte(sb.String()), 0600)
}

// quoteValue wraps the value in double quotes if it contains spaces,
// special shell characters, or is empty.
func quoteValue(v string) string {
	if v == "" || strings.ContainsAny(v, " \t\n\r#$&*(){}[]|;<>?`!\"'\\") {
		return fmt.Sprintf("%q", v)
	}
	return v
}
