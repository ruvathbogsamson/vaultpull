package vault

import (
	"strings"
	"testing"
)

func TestTemplateRenderer_Render_BasicSubstitution(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	r := NewTemplateRenderer(secrets)

	out, err := r.Render("host={{ index . \"DB_HOST\" }} port={{ index . \"DB_PORT\" }}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "host=localhost port=5432" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestTemplateRenderer_Render_MissingKey(t *testing.T) {
	r := NewTemplateRenderer(map[string]string{})

	_, err := r.Render("{{ index . \"MISSING\" }}")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestTemplateRenderer_Render_InvalidTemplate(t *testing.T) {
	r := NewTemplateRenderer(map[string]string{"K": "v"})

	_, err := r.Render("{{ .Unclosed")
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "template parse error") {
		t.Errorf("expected parse error message, got: %v", err)
	}
}

func TestTemplateRenderer_RenderAll_Success(t *testing.T) {
	secrets := map[string]string{"HOST": "db.local", "PORT": "3306"}
	r := NewTemplateRenderer(secrets)

	templates := map[string]string{
		"DSN": `{{ index . "HOST" }}:{{ index . "PORT" }}`,
		"URL": `mysql://{{ index . "HOST" }}`,
	}

	out, err := r.RenderAll(templates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "db.local:3306" {
		t.Errorf("DSN: got %q", out["DSN"])
	}
	if out["URL"] != "mysql://db.local" {
		t.Errorf("URL: got %q", out["URL"])
	}
}

func TestTemplateRenderer_RenderAll_PartialFailure(t *testing.T) {
	secrets := map[string]string{"HOST": "db.local"}
	r := NewTemplateRenderer(secrets)

	templates := map[string]string{
		"GOOD": `{{ index . "HOST" }}`,
		"BAD":  `{{ index . "MISSING_KEY" }}`,
	}

	_, err := r.RenderAll(templates)
	if err == nil {
		t.Fatal("expected error for missing key in RenderAll")
	}
}
