package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultpull/internal/vault"
)

// AuditDeviceFlags holds parsed flags for the audit-device command.
type AuditDeviceFlags struct {
	Address string
	Token   string
	Verbose bool
}

// ParseAuditDeviceFlags parses CLI flags for the audit-device subcommand.
func ParseAuditDeviceFlags(args []string) (*AuditDeviceFlags, error) {
	fs := flag.NewFlagSet("audit-device", flag.ContinueOnError)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault server address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	verbose := fs.Bool("verbose", false, "Enable verbose output")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *address == "" {
		return nil, fmt.Errorf("vault address is required (set -address or VAULT_ADDR)")
	}
	if *token == "" {
		return nil, fmt.Errorf("vault token is required (set -token or VAULT_TOKEN)")
	}
	return &AuditDeviceFlags{
		Address: *address,
		Token:   *token,
		Verbose: *verbose,
	}, nil
}

// RunAuditDevice executes the audit-device command.
func RunAuditDevice(flags *AuditDeviceFlags) error {
	c, err := vault.NewAuditDeviceClient(flags.Address, flags.Token)
	if err != nil {
		return fmt.Errorf("creating audit device client: %w", err)
	}

	devices, err := c.ListAuditDevices()
	if err != nil {
		return fmt.Errorf("listing audit devices: %w", err)
	}

	if len(devices) == 0 {
		fmt.Println("no audit devices enabled")
		return nil
	}

	for path, dev := range devices {
		fmt.Printf("path: %s  type: %s  description: %s\n", path, dev.Type, dev.Description)
		if flags.Verbose {
			for k, v := range dev.Options {
				fmt.Printf("  option %s = %s\n", k, v)
			}
		}
	}
	return nil
}
