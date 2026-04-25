package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/yourusername/vaultpull/internal/vault"
)

// PluginFlags holds parsed flags for the plugin subcommand.
type PluginFlags struct {
	Address    string
	Token      string
	PluginType string
}

// ParsePluginFlags parses plugin subcommand flags from args.
func ParsePluginFlags(args []string) (*PluginFlags, error) {
	fs := flag.NewFlagSet("plugin", flag.ContinueOnError)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault server address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	pluginType := fs.String("type", "", "Plugin type to list (auth, secret, database)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *address == "" {
		return nil, fmt.Errorf("vault address is required (set -address or VAULT_ADDR)")
	}
	if *token == "" {
		return nil, fmt.Errorf("vault token is required (set -token or VAULT_TOKEN)")
	}
	return &PluginFlags{
		Address:    *address,
		Token:      *token,
		PluginType: *pluginType,
	}, nil
}

// RunPlugin executes the plugin list command and writes results to w.
func RunPlugin(flags *PluginFlags, w io.Writer) error {
	client, err := vault.NewPluginClient(flags.Address, flags.Token)
	if err != nil {
		return fmt.Errorf("creating plugin client: %w", err)
	}

	plugins, err := client.ListPlugins(flags.PluginType)
	if err != nil {
		return fmt.Errorf("listing plugins: %w", err)
	}

	if len(plugins) == 0 {
		fmt.Fprintln(w, "no plugins found")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tTYPE\tVERSION\tBUILTIN")
	for _, p := range plugins {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%v\n", p.Name, p.Type, p.Version, p.Builtin)
	}
	return tw.Flush()
}
