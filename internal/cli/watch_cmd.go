package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

// WatchFlags holds parsed flags for the watch sub-command.
type WatchFlags struct {
	Interval  time.Duration
	OutputFile string
	Namespace  string
	Verbose    bool
}

// ParseWatchFlags parses watch sub-command flags from args.
func ParseWatchFlags(args []string, stderr io.Writer) (*WatchFlags, error) {
	fs := flag.NewFlagSet("watch", flag.ContinueOnError)
	fs.SetOutput(stderr)

	f := &WatchFlags{}
	fs.DurationVar(&f.Interval, "interval", 60*time.Second, "polling interval")
	fs.StringVar(&f.OutputFile, "output", ".env", "output .env file path")
	fs.StringVar(&f.Namespace, "namespace", "", "filter key namespace prefix")
	fs.BoolVar(&f.Verbose, "verbose", false, "enable verbose logging")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return f, nil
}

// RunWatch starts a watch loop that re-writes the .env file on secret changes.
func RunWatch(ctx context.Context, flags *WatchFlags, cfg *config.Config, out io.Writer) error {
	client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("watch: create client: %w", err)
	}

	onChange := func(path string, secrets map[string]string) {
		filtered := env.Filter(secrets, flags.Namespace)
		w, werr := env.NewWriter(flags.OutputFile)
		if werr != nil {
			fmt.Fprintf(out, "watch: open writer: %v\n", werr)
			return
		}
		if werr = w.Write(filtered); werr != nil {
			fmt.Fprintf(out, "watch: write env: %v\n", werr)
			return
		}
		if flags.Verbose {
			fmt.Fprintf(out, "watch: updateds (%d keys) from %s\n",
				flags.OutputFile, len(filtered), path)
		}
Error := func(path string, e error) {
		fmt.Fprintf(os.Stderr, "watch: error polling %s: %v\n", path, e)
	}

	watcher := vault.NewWatcher(client, cfg.SecretPath, vault.WatchConfig{
		Interval: flags.Interval,
		OnChange: onChange,
		OnError:  onError,
	})
	watcher.Start(ctx)
	defer watcher.Stop()

	<-ctx.Done()
	return nil
}
