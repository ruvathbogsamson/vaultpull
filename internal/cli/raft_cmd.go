package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/your-org/vaultpull/internal/vault"
)

// RaftFlags holds parsed flags for the raft subcommand.
type RaftFlags struct {
	Address string
	Token   string
}

// ParseRaftFlags parses raft subcommand flags from args.
func ParseRaftFlags(args []string) (*RaftFlags, error) {
	fs := flag.NewFlagSet("raft", flag.ContinueOnError)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return &RaftFlags{
		Address: *address,
		Token:   *token,
	}, nil
}

// RunRaft executes the raft configuration command, printing cluster peers to w.
func RunRaft(f *RaftFlags, w io.Writer) error {
	if f.Address == "" {
		return fmt.Errorf("vault address is required (set --address or VAULT_ADDR)")
	}
	if f.Token == "" {
		return fmt.Errorf("vault token is required (set --token or VAULT_TOKEN)")
	}

	client, err := vault.NewRaftClient(f.Address, f.Token)
	if err != nil {
		return fmt.Errorf("creating raft client: %w", err)
	}

	cfg, err := client.GetRaftConfig()
	if err != nil {
		return fmt.Errorf("fetching raft config: %w", err)
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NODE ID\tADDRESS\tLEADER\tVOTER")
	for _, s := range cfg.Servers {
		fmt.Fprintf(tw, "%s\t%s\t%v\t%v\n", s.NodeID, s.Address, s.Leader, s.Voter)
	}
	return tw.Flush()
}
