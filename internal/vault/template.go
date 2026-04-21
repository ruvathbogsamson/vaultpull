package vault

import (
	"bytes"
	"fmt"
	"text/template"
)

// TemplateRenderer renders secret values into Go text/template strings.
// It allows users to embed Vault secret values into arbitrary text output.
type TemplateRenderer struct {
	secrets map[string]string
}

// NewTemplateRenderer creates a TemplateRenderer backed by the provided secret map.
func NewTemplateRenderer(secrets map[string]string) *TemplateRenderer {
	return &TemplateRenderer{secrets: secrets}
}

// Render executes tmplStr as a Go template with the secrets map as its data.
// Secret keys are referenced as {{ .KEY_NAME }} in the template.
func (r *TemplateRenderer) Render(tmplStr string) (string, error) {
	tmpl, err := template.New("vault").Option("missingkey=error").Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, r.secrets); err != nil {
		return "", fmt.Errorf("template execute error: %w", err)
	}

	return buf.String(), nil
}

// RenderAll renders each value in templates map and returns the results.
// If any template fails, rendering stops and the error is returned.
func (r *TemplateRenderer) RenderAll(templates map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(templates))
	for key, tmplStr := range templates {
		result, err := r.Render(tmplStr)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", key, err)
		}
		out[key] = result
	}
	return out, nil
}
