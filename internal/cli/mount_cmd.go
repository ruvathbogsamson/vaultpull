package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/yourusername/vaultpull/internal/vault"
)

// MountFlags holds parsed flags for the mount subcommand.
type MountFlags struct {
	Address string
	Token   string
}

// ParseMountFlags parses mount subcommand flags from args.
func ParseMountFlags(args []string) (*MountFlags, error) {
	fs := flag.NewFlagSet("mount", flag.ContinueOnError)
	addr := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *addr == "" {
		return nil, fmt.Errorf("flag -address is required")
	}
	if *token == "" {
		return nil, fmt.Errorf("flag -token is required")
	}
	return &MountFlags{Address: *addr, Token: *token}, nil
}

// RunMount lists all Vault mount paths and writes them to out.
func RunMount(f *MountFlags, out io.Writer) error {
	c, err := vault.NewMountClient(f.Address, f.Token)
	if err != nil {
		return fmt.Errorf("initialising mount client: %w", err)
	}
	mounts, err := c.ListMountPaths()
	if err != nil {
		return fmt.Errorf("listing mounts: %w", err)
	}
	sort.Slice(mounts, func(i, j int) bool {
		return mounts[i].Path < mounts[j].Path
	})
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PATH\tTYPE\tDESCRIPTION")
	for _, m := range mounts {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", m.Path, m.Type, m.Description)
	}
	return tw.Flush()
}
