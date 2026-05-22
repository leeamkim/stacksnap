package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/stacksnap/internal/snapshot"
	"github.com/stacksnap/internal/validate"
)

func main() {
	snapshotFile := flag.String("snapshot", "", "path to snapshot JSON file (required)")
	requiredTools := flag.String("require", "", "comma-separated list of required tool names")
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

	var required []string
	if *requiredTools != "" {
		for _, r := range strings.Split(*requiredTools, ",") {
			if t := strings.TrimSpace(r); t != "" {
				required = append(required, t)
			}
		}
	}

	report, err := validate.Check(snap, required)
	if err != nil {
		fmt.Fprintf(os.Stderr, "validation error: %v\n", err)
		os.Exit(1)
	}

	for _, r := range report.Results {
		status := "OK  "
		if !r.Passed {
			status = "FAIL"
		}
		fmt.Printf("[%s] %s\n", status, r.Message)
	}

	fmt.Printf("\nResults: %d passed, %d failed\n", report.Passed, report.Failed)
	if report.Failed > 0 {
		os.Exit(1)
	}
}

func loadSnapshot(path string) (snapshot.Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return snapshot.Snapshot{}, err
	}
	var snap snapshot.Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return snapshot.Snapshot{}, err
	}
	return snap, nil
}
