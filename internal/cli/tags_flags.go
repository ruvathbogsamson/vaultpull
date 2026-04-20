package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// TagsFlags holds parsed CLI flags for tag-based secret filtering.
type TagsFlags struct {
	Tags   []string // raw "key=value" pairs
	DryRun bool
}

// ParseTagsFlags parses tag filter flags from args.
// Returns TagsFlags and any parse error.
func ParseTagsFlags(args []string) (*TagsFlags, error) {
	fs := flag.NewFlagSet("tags", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var rawTags string
	var dryRun bool

	fs.StringVar(&rawTags, "tags", "", "Comma-separated tag filters as key=value pairs (e.g. env=prod,team=backend)")
	fs.BoolVar(&dryRun, "dry-run", false, "Print matched secret keys without writing")

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("tags flags: %w", err)
	}

	var pairs []string
	if rawTags != "" {
		for _, p := range strings.Split(rawTags, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				pairs = append(pairs, p)
			}
		}
	}

	return &TagsFlags{
		Tags:   pairs,
		DryRun: dryRun,
	}, nil
}
