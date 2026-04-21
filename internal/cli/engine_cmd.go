package cli

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/example/vaultpull/internal/vault"
	vaultapi "github.com/hashicorp/vault/api"
)

// EngineFlags holds parsed flags for the engine subcommand.
type EngineFlags struct {
	Address string
	Token   string
	Mount   string
}

// ParseEngineFlags parses CLI flags for the engine subcommand.
func ParseEngineFlags(args []string) (*EngineFlags, error) {
	fs := flag.NewFlagSet("engine", flag.ContinueOnError)
	addr := fs.String("addr", "http://127.0.0.1:8200", "Vault server address")
	token := fs.String("token", "", "Vault token (or VAULT_TOKEN env)")
	mount := fs.String("mount", "", "Specific mount path to inspect (optional)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *token == "" {
		*token = os.Getenv("VAULT_TOKEN")
	}
	return &EngineFlags{Address: *addr, Token: *token, Mount: *mount}, nil
}

// RunEngine lists or inspects Vault secrets engines.
func RunEngine(flags *EngineFlags) error {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = flags.Address
	c, err := vaultapi.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}
	c.SetToken(flags.Token)

	ec := vault.NewEngineClient(c)

	if flags.Mount != "" {
		m, err := ec.GetMount(flags.Mount)
		if err != nil {
			return err
		}
		fmt.Printf("Path:        %s\nType:        %s\nDescription: %s\nVersion:     %s\n",
			m.Path, m.Type, m.Description, m.Version)
		return nil
	}

	mounts, err := ec.ListMounts()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATH\tTYPE\tVERSION\tDESCRIPTION")
	for _, m := range mounts {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", m.Path, m.Type, m.Version, m.Description)
	}
	return w.Flush()
}
