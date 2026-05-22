package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/stacksnap/internal/lint"
	"github.com/stacksnap/internal/snapshot"
)

func main() {
	snapshotFile := flag.String("snapshot", "", "path to snapshot JSON file (required)")
	jsonOutput := flag.Bool("json", false, "output issues as JSON")
	flag.Parse()

	if *snapshotFile == "" {
		fmt.Fprintln(os.Stderr, "error: --snapshot flag is required")
		flag.Usage()
		os.Exit(1)
	}

	snap, err := loadSnapshot(*snapshotFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading snapshot: %v\n", err)
		os.Exit(1)
	}

	result, err := lint.Check(snap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lint error: %v\n", err)
		os.Exit(1)
	}

	if *jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result.Issues); err != nil {
			fmt.Fprintf(os.Stderr, "json encode error: %v\n", err)
			os.Exit(1)
		}
	} else {
		if len(result.Issues) == 0 {
			fmt.Println("✓ snapshot looks good — no issues found")
		} else {
			for _, iss := range result.Issues {
				fmt.Println(iss)
			}
		}
	}

	if !result.OK() {
		os.Exit(2)
	}
}

func loadSnapshot(path string) (*snapshot.Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var snap snapshot.Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return &snap, nil
}
