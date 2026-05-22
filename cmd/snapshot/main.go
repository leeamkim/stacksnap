package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/stacksnap/internal/export"
	"github.com/stacksnap/internal/snapshot"
)

func main() {
	var (
		format   string
		outDir   string
		filename string
		allowEnv string
	)

	flag.StringVar(&format, "format", "json", "Output format (json)")
	flag.StringVar(&outDir, "out", ".", "Output directory for the snapshot file")
	flag.StringVar(&filename, "filename", "", "Override output filename (default: auto-generated)")
	flag.StringVar(&allowEnv, "allow-env", "", "Comma-separated list of additional env var prefixes to include")
	flag.Parse()

	extraEnv := []string{}
	if allowEnv != "" {
		for _, prefix := range strings.Split(allowEnv, ",") {
			prefix = strings.TrimSpace(prefix)
			if prefix != "" {
				extraEnv = append(extraEnv, prefix)
			}
		}
	}

	fmt.Println("Capturing dev stack snapshot...")

	snap, err := snapshot.Capture(extraEnv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error capturing snapshot: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Detected %d tools, %d env vars\n", len(snap.Tools), len(snap.Env))

	opts := export.DefaultOptions()
	opts.Format = format
	opts.OutputDir = outDir
	if filename != "" {
		opts.Filename = filename
	}

	outPath, err := export.Export(snap, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error exporting snapshot: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Snapshot written to: %s\n", outPath)
}
