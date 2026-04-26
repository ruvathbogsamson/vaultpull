package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/vaultpull/internal/vault"
)

// MFAFlags holds parsed flags for the mfa subcommand.
type MFAFlags struct {
	Address   string
	Token     string
	RequestID string
	Payload   string
}

// ParseMFAFlags parses command-line flags for the mfa subcommand.
func ParseMFAFlags(args []string) (*MFAFlags, error) {
	fs := flag.NewFlagSet("mfa", flag.ContinueOnError)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	requestID := fs.String("request-id", "", "MFA request ID (required)")
	payload := fs.String("payload", "", "MFA payload value (e.g. TOTP code)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *requestID == "" {
		return nil, fmt.Errorf("mfa: --request-id is required")
	}
	return &MFAFlags{
		Address:   *address,
		Token:     *token,
		RequestID: *requestID,
		Payload:   *payload,
	}, nil
}

// RunMFA executes the mfa validation flow using the provided flags.
func RunMFA(f *MFAFlags) error {
	client, err := vault.NewMFAClient(f.Address, f.Token)
	if err != nil {
		return fmt.Errorf("mfa: failed to create client: %w", err)
	}

	req := vault.MFAValidateRequest{
		MFARequestID: f.RequestID,
		MFAPayload:   map[string]string{"totp": f.Payload},
	}

	resp, err := client.Validate(req)
	if err != nil {
		return fmt.Errorf("mfa: validation failed: %w", err)
	}

	fmt.Printf("MFA validated. Client token: %s\n", resp.Token)
	if len(resp.Policies) > 0 {
		fmt.Printf("Policies: %v\n", resp.Policies)
	}
	return nil
}
