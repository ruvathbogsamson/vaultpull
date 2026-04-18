package cli

import (
	"flag"
	"os"
)

// Options holds parsed CLI flag values.
type Options struct {
	OutputFile string
	Namespace  string
	DryRun     bool
	Verbose    bool
}

// ParseFlags parses command-line flags and returns Options.
func ParseFlags(args []string) (*Options, error) {
	fs := flag.NewFlagSet("vaultpull", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	opts := &Options{}

	fs.StringVar(&opts.OutputFile, "output", ".env", "path to output .env file")
	fs.StringVar(&opts.Namespace, "namespace", "", "filter secrets by namespace prefix")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "print secrets without writing to file")
	fs.BoolVar(&opts.Verbose, "verbose", false, "enable verbose logging")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	return opts, nil
}
