package config

import (
	"os"
	"testing"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	t.Setenv(key, value)
}

func TestLoad_MissingToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")
	os.Unsetenv("VAULTPULL_VAULT_TOKEN")

	_, err := Load("")
	if err == nil {
		t.Fatal("expected error when vault_token is missing, got nil")
	}
}

func TestLoad_MissingSecretPath(t *testing.T) {
	setEnv(t, "VAULT_TOKEN", "root")
	os.Unsetenv("VAULTPULL_SECRET_PATH")

	_, err := Load("")
	if err == nil {
		t.Fatal("expected error when secret_path is missing, got nil")
	}
}

func TestLoad_Defaults(t *testing.T) {
	setEnv(t, "VAULT_TOKEN", "root")
	setEnv(t, "VAULTPULL_SECRET_PATH", "secret/data/myapp")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("expected default vault_addr, got %q", cfg.VaultAddr)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default output_file '.env', got %q", cfg.OutputFile)
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	setEnv(t, "VAULT_TOKEN", "s.mytoken")
	setEnv(t, "VAULT_ADDR", "https://vault.example.com")
	setEnv(t, "VAULTPULL_SECRET_PATH", "kv/data/service")
	setEnv(t, "VAULTPULL_NAMESPACE", "engineering")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Namespace != "engineering" {
		t.Errorf("expected namespace 'engineering', got %q", cfg.Namespace)
	}
	if cfg.SecretPath != "kv/data/service" {
		t.Errorf("expected secret_path 'kv/data/service', got %q", cfg.SecretPath)
	}
}
