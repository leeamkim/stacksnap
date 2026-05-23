// Command compare checks a snapshot against a named template baseline.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/stacksnap/internal/compare"
	"github.com/stacksnap/internal/snapshot"
	"github.com/stacksnap/internal/template"
)

func main() {
	snapshotPath := flag.String("snapshot", "", "path to snapshot JSON file (required)")
	templatePath := flag.String("template", "", "path to template JSON file (required)")
	flag.Parse()

	if *snapshotPath == "" || *templatePath == "" {
		fmt.Fprintln(os.Stderr, "error: --snapshot and --template are required")
		flag.Usage()
		os.Exit(1)
	}

	snap, err := loadSnapshot(*snapshotPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading snapshot: %v\n", err)
		os.Exit(1)
	}

	tmpl, err := template.Load(*templatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading template: %v\n", err)
		os.Exit(1)
	}

	res, err := compare.AgainstTemplate(snap, tmpl, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error comparing: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(res.Summary())

	if len(res.Missing) > 0 {
		fmt.Println("missing:", res.Missing)
	}
	if len(res.Mismatch) > 0 {
		fmt.Println("mismatch:", res.Mismatch)
	}
	if len(res.Extra) > 0 {
		fmt.Println("extra:", res.Extra)
	}

	if !res.OK() {
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
		return nil, err
	}
	return &snap, nil
}
