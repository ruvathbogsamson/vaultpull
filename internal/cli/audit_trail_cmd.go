package cli

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/your-org/vaultpull/internal/vault"
)

// AuditTrailFlags holds parsed flags for the audit-trail command.
type AuditTrailFlags struct {
	Path      string
	Namespace string
	Verbose   bool
}

// ParseAuditTrailFlags parses CLI flags for the audit-trail subcommand.
func ParseAuditTrailFlags(args []string) (*AuditTrailFlags, error) {
	fs := flag.NewFlagSet("audit-trail", flag.ContinueOnError)
	path := fs.String("path", "", "Vault secret path to audit")
	namespace := fs.String("namespace", "", "Filter events by namespace prefix")
	verbose := fs.Bool("verbose", false, "Show full event details")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *path == "" {
		return nil, fmt.Errorf("flag -path is required")
	}
	return &AuditTrailFlags{
		Path:      *path,
		Namespace: *namespace,
		Verbose:   *verbose,
	}, nil
}

// RunAuditTrail executes the audit-trail command using the provided trail.
func RunAuditTrail(flags *AuditTrailFlags, trail *vault.AuditTrail) error {
	events := trail.Events()
	if len(events) == 0 {
		fmt.Fprintln(os.Stdout, "No audit events found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tOPERATION\tPATH\tNAMESPACE\tERROR")

	for _, ev := range events {
		if flags.Namespace != "" && ev.Namespace != flags.Namespace {
			continue
		}
		errStr := "-"
		if ev.Error != "" {
			errStr = ev.Error
		}
		ns := ev.Namespace
		if ns == "" {
			ns = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			ev.Timestamp.Format(time.RFC3339),
			ev.Operation,
			ev.Path,
			ns,
			errStr,
		)
	}
	return w.Flush()
}
