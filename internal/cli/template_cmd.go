package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/vault"
)

// TemplateCmdFlags holds parsed flags for the template subcommand.
type TemplateCmdFlags struct {
	TemplatePath string
	OutputPath   string
	SecretPath   string
	Verbose      bool
}

// ParseTemplateFlags parses template subcommand flags from args.
func ParseTemplateFlags(args []string) (*TemplateCmdFlags, error) {
	fs := flag.NewFlagSet("template", flag.ContinueOnError)

	tmplPath := fs.String("template", "", "Path to the template file (required)")
	outPath := fs.String("output", ".env", "Path to write rendered output")
	secretPath := fs.String("secret-path", "", "Vault secret path to use as template data")
	verbose := fs.Bool("verbose", false, "Enable verbose output")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *tmplPath == "" {
		return nil, fmt.Errorf("--template is required")
	}
	return &TemplateCmdFlags{
		TemplatePath: *tmplPath,
		OutputPath:   *outPath,
		SecretPath:   *secretPath,
		Verbose:      *verbose,
	}, nil
}

// RunTemplate executes the template rendering pipeline.
func RunTemplate(flags *TemplateCmdFlags) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	if flags.SecretPath != "" {
		cfg.SecretPath = flags.SecretPath
	}

	client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := vault.FetchSecrets(client, cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	tmplBytes, err := os.ReadFile(flags.TemplatePath)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}

	renderer := vault.NewTemplateRenderer(secrets)
	output, err := renderer.Render(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	if flags.Verbose {
		fmt.Printf("Writing rendered template to %s\n", flags.OutputPath)
	}
	return os.WriteFile(flags.OutputPath, []byte(output), 0600)
}
