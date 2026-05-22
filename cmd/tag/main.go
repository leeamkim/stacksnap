package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/stacksnap/internal/snapshot"
	"github.com/stacksnap/internal/tag"
)

func main() {
	snapshotFile := flag.String("snapshot", "", "path to snapshot JSON file (required)")
	addFlag := flag.String("add", "", "comma-separated tags to add")
	removeFlag := flag.String("remove", "", "comma-separated tags to remove")
	listFlag := flag.Bool("list", false, "list current tags and exit")
	flag.Parse()

	if *snapshotFile == "" {
		fmt.Fprintln(os.Stderr, "error: --snapshot is required")
		os.Exit(1)
	}

	snap, err := loadSnapshot(*snapshotFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading snapshot: %v\n", err)
		os.Exit(1)
	}

	if *listFlag {
		for _, t := range tag.List(snap) {
			fmt.Println(t)
		}
		return
	}

	if *addFlag != "" {
		parts := splitCSV(*addFlag)
		if err := tag.Add(snap, parts...); err != nil {
			fmt.Fprintf(os.Stderr, "error adding tags: %v\n", err)
			os.Exit(1)
		}
	}

	if *removeFlag != "" {
		tag.Remove(snap, splitCSV(*removeFlag)...)
	}

	if err := saveSnapshot(*snapshotFile, snap); err != nil {
		fmt.Fprintf(os.Stderr, "error saving snapshot: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("tags updated: %v\n", tag.List(snap))
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func loadSnapshot(path string) (*snapshot.Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s snapshot.Snapshot
	return &s, json.Unmarshal(data, &s)
}

func saveSnapshot(path string, s *snapshot.Snapshot) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
