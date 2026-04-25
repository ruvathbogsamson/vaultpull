package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/vault"
)

// TOTPFlags holds parsed flags for the totp subcommand.
type TOTPFlags struct {
	Address string
	Token   string
	Mount   string
	KeyName string
	Code    string
	Action  string // "generate" or "validate"
}

// ParseTOTPFlags parses CLI flags for the totp subcommand.
func ParseTOTPFlags(args []string) (*TOTPFlags, error) {
	fs := flag.NewFlagSet("totp", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault server address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	mount := fs.String("mount", "totp", "TOTP secrets engine mount path")
	keyName := fs.String("key", "", "TOTP key name (required)")
	code := fs.String("code", "", "TOTP code to validate (required for validate action)")
	action := fs.String("action", "generate", "Action to perform: generate or validate")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if *keyName == "" {
		return nil, fmt.Errorf("flag -key is required")
	}
	if *action != "generate" && *action != "validate" {
		return nil, fmt.Errorf("invalid action %q: must be generate or validate", *action)
	}
	if *action == "validate" && *code == "" {
		return nil, fmt.Errorf("flag -code is required for validate action")
	}

	return &TOTPFlags{
		Address: *address,
		Token:   *token,
		Mount:   *mount,
		KeyName: *keyName,
		Code:    *code,
		Action:  *action,
	}, nil
}

// RunTOTP executes the totp subcommand.
func RunTOTP(f *TOTPFlags) error {
	client, err := vault.NewTOTPClient(f.Address, f.Token, f.Mount)
	if err != nil {
		return fmt.Errorf("totp client: %w", err)
	}

	switch f.Action {
	case "generate":
		code, err := client.GenerateCode(f.KeyName)
		if err != nil {
			return fmt.Errorf("generate code: %w", err)
		}
		fmt.Println(code)
	case "validate":
		valid, err := client.ValidateCode(f.KeyName, f.Code)
		if err != nil {
			return fmt.Errorf("validate code: %w", err)
		}
		if valid {
			fmt.Println("valid")
		} else {
			fmt.Println("invalid")
			return fmt.Errorf("code is not valid")
		}
	}
	return nil
}
