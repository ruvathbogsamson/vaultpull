package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/example/vaultpull/internal/vault"
)

// AWSFlags holds parsed flags for the AWS credentials command.
type AWSFlags struct {
	Address string
	Token   string
	Mount   string
	Role    string
}

// ParseAWSFlags parses AWS subcommand flags from args.
func ParseAWSFlags(args []string) (*AWSFlags, error) {
	fs := flag.NewFlagSet("aws", flag.ContinueOnError)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault server address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	mount := fs.String("mount", "aws", "AWS secrets engine mount path")
	role := fs.String("role", "", "AWS role name to generate credentials for (required)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *role == "" {
		return nil, fmt.Errorf("flag -role is required")
	}
	return &AWSFlags{
		Address: *address,
		Token:   *token,
		Mount:   *mount,
		Role:    *role,
	}, nil
}

// RunAWS executes the AWS credentials command, writing results to out.
func RunAWS(args []string, out io.Writer) error {
	flags, err := ParseAWSFlags(args)
	if err != nil {
		return err
	}

	client, err := vault.NewAWSClient(flags.Address, flags.Token, flags.Mount)
	if err != nil {
		return fmt.Errorf("creating aws client: %w", err)
	}

	creds, err := client.GenerateCredentials(flags.Role)
	if err != nil {
		return fmt.Errorf("generating credentials: %w", err)
	}

	fmt.Fprintf(out, "AWS_ACCESS_KEY_ID=%s\n", creds.AccessKey)
	fmt.Fprintf(out, "AWS_SECRET_ACCESS_KEY=%s\n", creds.SecretKey)
	if creds.SecurityToken != "" {
		fmt.Fprintf(out, "AWS_SESSION_TOKEN=%s\n", creds.SecurityToken)
	}
	fmt.Fprintf(out, "# lease_duration: %ds\n", creds.LeaseDuration)
	return nil
}
