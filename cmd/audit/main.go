// Command audit prints the audit trail of a snapshot file.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/stacksnap/internal/audit"
	"github.com/stacksnap/internal/snapshot"
)

func main() {
	snapshotPath := flag.String("snapshot", "", "path to snapshot JSON file (required)")
	actor := flag.String("actor", "", "actor name to record an access event (optional)")
	flag.Parse()

	if *snapshotPath == "" {
		fmt.Fprintln(os.Stderr, "error: -snapshot flag is required")
		flag.Usage()
		os.Exit(1)
	}

	snap, err := loadSnapshot(*snapshotPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading snapshot: %v\n", err)
		os.Exit(1)
	}

	if *actor != "" {
		if err := audit.Record(snap, audit.EventShare, *actor, "accessed via audit command"); err != nil {
			fmt.Fprintf(os.Stderr, "error recording audit event: %v\n", err)
			os.Exit(1)
		}
	}

	events, err := audit.List(snap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing audit events: %v\n", err)
		os.Exit(1)
	}

	if len(events) == 0 {
		fmt.Println("no audit events recorded")
		return
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(events); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding events: %v\n", err)
		os.Exit(1)
	}
}

func loadSnapshot(path string) (*snapshot.Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var snap snapshot.Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}
