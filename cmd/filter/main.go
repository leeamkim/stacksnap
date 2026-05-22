package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/stacksnap/internal/filter"
	"github.com/stacksnap/internal/snapshot"
)

func main() {
	snapshotFile := flag.String("snapshot", "", "path to snapshot JSON file (required)")
	noTools := flag.Bool("no-tools", false, "exclude tools from output")
	noEnv := flag.Bool("no-env", false, "exclude env vars from output")
	noOS := flag.Bool("no-os", false, "exclude OS info from output")
	tagsRaw := flag.String("tags", "", "comma-separated list of tool name substrings to include")
	outFile := flag.String("out", "", "write output to file instead of stdout")
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

	opts := filter.DefaultOptions()
	if *noTools {
		opts.Tools = false
	}
	if *noEnv {
		opts.Env = false
	}
	if *noOS {
		opts.OS = false
	}
	if *tagsRaw != "" {
		opts.Tags = strings.Split(*tagsRaw, ",")
	}

	result := filter.Apply(snap, opts)

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encoding result: %v\n", err)
		os.Exit(1)
	}

	if *outFile != "" {
		if err := os.WriteFile(*outFile, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "filtered snapshot written to %s\n", *outFile)
		return
	}

	fmt.Println(string(data))
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
