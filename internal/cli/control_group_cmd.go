package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/example/vaultpull/internal/vault"
)

// ControlGroupFlags holds parsed flags for the control-group subcommand.
type ControlGroupFlags struct {
	Address  string
	Token    string
	Accessor string
	Action   string // "authorize" or "check"
}

// ParseControlGroupFlags parses command-line flags for the control-group command.
func ParseControlGroupFlags(args []string) (ControlGroupFlags, error) {
	fs := flag.NewFlagSet("control-group", flag.ContinueOnError)

	var f ControlGroupFlags
	fs.StringVar(&f.Address, "address", os.Getenv("VAULT_ADDR"), "Vault server address")
	fs.StringVar(&f.Token, "token", os.Getenv("VAULT_TOKEN"), "Vault token")
	fs.StringVar(&f.Accessor, "accessor", "", "Control group accessor token")
	fs.StringVar(&f.Action, "action", "check", "Action to perform: authorize or check")

	if err := fs.Parse(args); err != nil {
		return ControlGroupFlags{}, err
	}

	if f.Accessor == "" {
		return ControlGroupFlags{}, fmt.Errorf("--accessor is required")
	}

	if f.Action != "authorize" && f.Action != "check" {
		return ControlGroupFlags{}, fmt.Errorf("--action must be \"authorize\" or \"check\", got %q", f.Action)
	}

	return f, nil
}

// RunControlGroup executes the control-group command using parsed flags.
func RunControlGroup(f ControlGroupFlags) error {
	client, err := vault.NewControlGroupClient(f.Address, f.Token)
	if err != nil {
		return fmt.Errorf("control-group: %w", err)
	}

	switch f.Action {
	case "authorize":
		ok, err := client.Authorize(f.Accessor)
		if err != nil {
			return fmt.Errorf("authorize failed: %w", err)
		}
		if ok {
			fmt.Println("Authorization successful.")
		} else {
			fmt.Println("Authorization denied.")
		}

	case "check":
		status, err := client.CheckRequest(f.Accessor)
		if err != nil {
			return fmt.Errorf("check request failed: %w", err)
		}
		fmt.Printf("Request ID : %s\n", status.RequestID)
		fmt.Printf("Accessor   : %s\n", status.Accessor)
		fmt.Printf("Approved   : %v\n", status.Approved)
	}

	return nil
}
