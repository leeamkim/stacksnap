package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/stacksnap/internal/restore"
)

func main() {
	var (
		snapshotFile = flag.String("snapshot", "", "path to snapshot JSON file (required)")
		dryRun       = flag.Bool("dry-run", false, "print commands without executing them")
		verbose      = flag.Bool("verbose", false, "print each action")
	)
	flag.Parse()

	if *snapshotFile == "" {
		fmt.Fprintln(os.Stderr, "error: --snapshot flag is required")
		flag.Usage()
		os.Exit(1)
	}

	snap, err := restore.LoadSnapshot(*snapshotFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading snapshot: %v\n", err)
		os.Exit(1)
	}

	opts := restore.Options{
		DryRun:  *dryRun,
		Verbose: *verbose,
	}

	if *dryRun {
		fmt.Println("[dry-run] the following commands would be executed:")
	}

	if err := restore.Apply(snap, opts); err != nil {
		fmt.Fprintf(os.Stderr, "restore failed: %v\n", err)
		os.Exit(1)
	}

	if !*dryRun {
		fmt.Printf("restore complete: applied %d tool(s) from %s\n", len(snap.Tools), *snapshotFile)
	}
}
