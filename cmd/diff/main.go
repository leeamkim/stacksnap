package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/stacksnap/internal/diff"
	"github.com/stacksnap/internal/snapshot"
)

func main() {
	basePath := flag.String("base", "", "path to baseline snapshot JSON file (required)")
	currPath := flag.String("current", "", "path to current snapshot JSON file (required)")
	flag.Parse()

	if *basePath == "" || *currPath == "" {
		fmt.Fprintln(os.Stderr, "usage: stacksnap-diff -base <file> -current <file>")
		os.Exit(1)
	}

	baseline, err := loadSnapshot(*basePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading baseline: %v\n", err)
		os.Exit(1)
	}

	current, err := loadSnapshot(*currPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading current snapshot: %v\n", err)
		os.Exit(1)
	}

	result := diff.Compare(baseline, current)
	fmt.Print(diff.Summary(result))

	if !result.Equal {
		os.Exit(2)
	}
}

func loadSnapshot(path string) (*snapshot.Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer f.Close()

	var s snapshot.Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("decode %q: %w", path, err)
	}
	return &s, nil
}
