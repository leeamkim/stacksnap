package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/stacksnap/internal/pin"
	"github.com/stacksnap/internal/snapshot"
)

func main() {
	snapshotFile := flag.String("snapshot", "", "path to snapshot JSON file (required)")
	tool := flag.String("tool", "", "tool name to pin or unpin")
	version := flag.String("version", "", "version to pin")
	author := flag.String("author", "", "who is pinning (optional)")
	unpinFlag := flag.Bool("unpin", false, "remove pin for the given tool")
	getFlag := flag.Bool("get", false, "retrieve pin info for the given tool")
	flag.Parse()

	if *snapshotFile == "" {
		fmt.Fprintln(os.Stderr, "error: --snapshot is required")
		os.Exit(1)
	}
	if *tool == "" {
		fmt.Fprintln(os.Stderr, "error: --tool is required")
		os.Exit(1)
	}

	snap, err := loadSnapshot(*snapshotFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading snapshot: %v\n", err)
		os.Exit(1)
	}

	switch {
	case *getFlag:
		p, err := pin.GetPin(snap, *tool)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if p == nil {
			fmt.Printf("no pin found for tool %q\n", *tool)
		} else {
			fmt.Printf("tool: %s\nversion: %s\npinned_by: %s\npinned_at: %s\n",
				p.Name, p.Version, p.PinnedBy, p.PinnedAt.Format("2006-01-02T15:04:05Z"))
		}

	case *unpinFlag:
		if err := pin.Unpin(snap, *tool); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if err := saveSnapshot(*snapshotFile, snap); err != nil {
			fmt.Fprintf(os.Stderr, "error saving snapshot: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("unpinned tool %q\n", *tool)

	default:
		if *version == "" {
			fmt.Fprintln(os.Stderr, "error: --version is required for pinning")
			os.Exit(1)
		}
		if err := pin.Pin(snap, *tool, *version, *author); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if err := saveSnapshot(*snapshotFile, snap); err != nil {
			fmt.Fprintf(os.Stderr, "error saving snapshot: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("pinned %s@%s\n", *tool, *version)
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

func saveSnapshot(path string, snap *snapshot.Snapshot) error {
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
