package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/vault"
)

// TransitFlags holds parsed flags for transit encrypt/decrypt commands.
type TransitFlags struct {
	Address  string
	Token    string
	Mount    string
	KeyName  string
	Action   string // "encrypt" or "decrypt"
	Payload  string
}

// ParseTransitFlags parses CLI flags for the transit subcommand.
func ParseTransitFlags(args []string) (*TransitFlags, error) {
	fs := flag.NewFlagSet("transit", flag.ContinueOnError)
	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	mount := fs.String("mount", "transit", "Transit engine mount path")
	key := fs.String("key", "", "Transit key name (required)")
	action := fs.String("action", "encrypt", "Action: encrypt or decrypt")
	payload := fs.String("payload", "", "Plaintext (encrypt) or ciphertext (decrypt) (required)")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *key == "" {
		return nil, fmt.Errorf("transit: -key is required")
	}
	if *payload == "" {
		return nil, fmt.Errorf("transit: -payload is required")
	}
	if *action != "encrypt" && *action != "decrypt" {
		return nil, fmt.Errorf("transit: -action must be 'encrypt' or 'decrypt', got %q", *action)
	}
	return &TransitFlags{
		Address: *address,
		Token:   *token,
		Mount:   *mount,
		KeyName: *key,
		Action:  *action,
		Payload: *payload,
	}, nil
}

// RunTransit executes the transit encrypt or decrypt operation.
func RunTransit(f *TransitFlags) error {
	client, err := vault.NewTransitClient(f.Address, f.Token, f.Mount)
	if err != nil {
		return fmt.Errorf("transit: failed to create client: %w", err)
	}
	switch f.Action {
	case "encrypt":
		result, err := client.Encrypt(f.KeyName, f.Payload)
		if err != nil {
			return err
		}
		fmt.Println(result)
	case "decrypt":
		result, err := client.Decrypt(f.KeyName, f.Payload)
		if err != nil {
			return err
		}
		fmt.Println(result)
	}
	return nil
}
