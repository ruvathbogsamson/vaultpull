package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/user/vaultpull/internal/vault"
)

// AuthMethodFlags holds parsed flags for the auth-method command.
type AuthMethodFlags struct {
	Address string
	Token   string
}

// ParseAuthMethodFlags parses CLI flags for the auth-method subcommand.
func ParseAuthMethodFlags(args []string) (*AuthMethodFlags, error) {
	fs := flag.NewFlagSet("auth-method", flag.ContinueOnError)
	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault server address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return &AuthMethodFlags{Address: *address, Token: *token}, nil
}

// RunAuthMethod lists all enabled auth methods and prints them in a table.
func RunAuthMethod(f *AuthMethodFlags, out io.Writer) error {
	c, err := vault.NewAuthMethodClient(f.Address, f.Token)
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}

	methods, err := c.ListAuthMethods()
	if err != nil {
		return fmt.Errorf("listing auth methods: %w", err)
	}

	paths := make([]string, 0, len(methods))
	for p := range methods {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATH\tTYPE\tDESCRIPTION\tLOCAL")
	for _, p := range paths {
		m := methods[p]
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", p, m.Type, m.Description, m.Local)
	}
	return w.Flush()
}
