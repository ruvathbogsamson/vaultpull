package cli

import (
	"testing"
)

func TestParseKubernetesFlags_Defaults(t *testing.T) {
	flags, err := ParseKubernetesFlags([]string{
		"-role", "my-role",
		"-jwt", "my-jwt",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.Mount != "kubernetes" {
		t.Errorf("expected default mount 'kubernetes', got %q", flags.Mount)
	}
	if flags.Role != "my-role" {
		t.Errorf("expected role 'my-role', got %q", flags.Role)
	}
}

func TestParseKubernetesFlags_AllFlags(t *testing.T) {
	flags, err := ParseKubernetesFlags([]string{
		"-address", "http://vault:8200",
		"-mount", "k8s",
		"-role", "dev-role",
		"-jwt", "eyJhbGciOiJSUzI1NiJ9.test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.Address != "http://vault:8200" {
		t.Errorf("expected address 'http://vault:8200', got %q", flags.Address)
	}
	if flags.Mount != "k8s" {
		t.Errorf("expected mount 'k8s', got %q", flags.Mount)
	}
	if flags.JWT != "eyJhbGciOiJSUzI1NiJ9.test" {
		t.Errorf("unexpected jwt value: %q", flags.JWT)
	}
}

func TestParseKubernetesFlags_MissingRole(t *testing.T) {
	_, err := ParseKubernetesFlags([]string{"-jwt", "some-jwt"})
	if err == nil {
		t.Fatal("expected error for missing role")
	}
}

func TestParseKubernetesFlags_MissingJWT(t *testing.T) {
	_, err := ParseKubernetesFlags([]string{"-role", "my-role"})
	if err == nil {
		t.Fatal("expected error for missing jwt")
	}
}

func TestParseKubernetesFlags_InvalidFlag(t *testing.T) {
	_, err := ParseKubernetesFlags([]string{"-unknown", "value"})
	if err == nil {
		t.Fatal("expected error for invalid flag")
	}
}
