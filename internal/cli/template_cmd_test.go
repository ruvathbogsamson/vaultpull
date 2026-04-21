package cli

import (
	"strings"
	"testing"
)

func TestParseTemplateFlags_Defaults(t *testing.T) {
	flags, err := ParseTemplateFlags([]string{"--template", "tmpl.txt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.TemplatePath != "tmpl.txt" {
		t.Errorf("TemplatePath: got %q, want %q", flags.TemplatePath, "tmpl.txt")
	}
	if flags.OutputPath != ".env" {
		t.Errorf("OutputPath: got %q, want %q", flags.OutputPath, ".env")
	}
	if flags.Verbose {
		t.Error("Verbose should default to false")
	}
}

func TestParseTemplateFlags_AllFlags(t *testing.T) {
	flags, err := ParseTemplateFlags([]string{
		"--template", "my.tmpl",
		"--output", "out.env",
		"--secret-path", "secret/data/app",
		"--verbose",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.TemplatePath != "my.tmpl" {
		t.Errorf("TemplatePath: got %q", flags.TemplatePath)
	}
	if flags.OutputPath != "out.env" {
		t.Errorf("OutputPath: got %q", flags.OutputPath)
	}
	if flags.SecretPath != "secret/data/app" {
		t.Errorf("SecretPath: got %q", flags.SecretPath)
	}
	if !flags.Verbose {
		t.Error("expected Verbose=true")
	}
}

func TestParseTemplateFlags_MissingTemplate(t *testing.T) {
	_, err := ParseTemplateFlags([]string{})
	if err == nil {
		t.Fatal("expected error for missing --template flag")
	}
	if !strings.Contains(err.Error(), "--template") {
		t.Errorf("expected --template in error, got: %v", err)
	}
}

func TestParseTemplateFlags_InvalidFlag(t *testing.T) {
	_, err := ParseTemplateFlags([]string{"--template", "t.tmpl", "--unknown", "val"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}
