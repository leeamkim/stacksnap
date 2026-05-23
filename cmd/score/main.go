package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/stacksnap/internal/score"
	"github.com/stacksnap/internal/snapshot"
)

func main() {
	snapshotFile := flag.String("snapshot", "", "path to snapshot JSON file (required)")
	jsonOut := flag.Bool("json", false, "output result as JSON")
	flag.Parse()

	if *snapshotFile == "" {
		fmt.Fprintln(os.Stderr, "error: --snapshot flag is required")
		os.Exit(1)
	}

	snap, err := loadSnapshot(*snapshotFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading snapshot: %v\n", err)
		os.Exit(1)
	}

	res, err := score.Evaluate(snap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error evaluating snapshot: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(res); err != nil {
			fmt.Fprintf(os.Stderr, "error encoding result: %v\n", err)
			os.Exit(1)
		}
		return
	}

	fmt.Printf("Score : %d / 100\n", res.Score)
	fmt.Printf("Grade : %s\n", res.Grade)
	if len(res.Reasons) > 0 {
		fmt.Println("Issues:")
		for _, r := range res.Reasons {
			fmt.Printf("  - %s\n", r)
		}
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
