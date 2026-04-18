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
func (w *Writer) Write(secrets map[string]string) error {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, secrets[k]))
	}

	return os.WriteFile(w.filePath, []byte(sb.String()), 0600)
}
