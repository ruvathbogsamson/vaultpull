# vaultpull

A CLI tool to sync HashiCorp Vault secrets to local `.env` files with namespace filtering.

## Features

- Pulls secrets from a Vault KV v2 mount into a local `.env` file
- Filters secrets by namespace prefix (e.g. `APP_`, `DB_`)
- Dry-run mode to preview changes without writing to disk
- Verbose logging for debugging
- Configurable via environment variables or CLI flags

## Installation

```bash
go install github.com/yourorg/vaultpull/cmd/vaultpull@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/vaultpull.git
cd vaultpull
go build -o vaultpull ./cmd/vaultpull
```

## Usage

```
vaultpull [flags]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-output` | `.env` | Path to the output `.env` file |
| `-namespace` | `` | Filter secrets by key prefix (e.g. `APP`) |
| `-dry-run` | `false` | Preview output without writing to disk |
| `-verbose` | `false` | Enable verbose logging |

### Environment Variables

| Variable | Description |
|----------|-------------|
| `VAULT_ADDR` | Vault server address (default: `http://127.0.0.1:8200`) |
| `VAULT_TOKEN` | Vault authentication token (**required**) |
| `VAULT_SECRET_PATH` | Path to the secret in Vault (**required**, e.g. `secret/myapp`) |
| `VAULT_NAMESPACE` | Filter prefix — overridden by `-namespace` flag if set |

## Example

```bash
export VAULT_ADDR=https://vault.example.com
export VAULT_TOKEN=s.mytoken
export VAULT_SECRET_PATH=secret/data/myapp

# Write all secrets to .env
vaultpull -output .env

# Preview only APP_ prefixed secrets
vaultpull -namespace APP -dry-run -verbose
```

Example output `.env`:

```dotenv
APP_DATABASE_URL=postgres://localhost/mydb
APP_SECRET_KEY=supersecret
```

## Configuration Precedence

1. CLI flags (highest priority)
2. Environment variables
3. Built-in defaults (lowest priority)

## Development

```bash
# Run all tests
go test ./...

# Run with race detector
go test -race ./...

# Build
go build ./cmd/vaultpull
```

## License

MIT
