// Command annotate adds or shows annotations on a snapshot file.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/stacksnap/internal/annotate"
	"github.com/stacksnap/internal/snapshot"
)

func main() {
	snapshotFile := flag.String("snapshot", "", "path to snapshot JSON file (required)")
	author := flag.String("author", "", "author name for the annotation")
	note := flag.String("note", "", "annotation text to add")
	show := flag.Bool("show", false, "print existing annotation and exit")
	clear := flag.Bool("clear", false, "remove existing annotation")
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

	switch {
	case *show:
		a := annotate.Get(snap)
		if a == nil {
			fmt.Println("no annotation found")
			return
		}
		printAnnotation(a)

	case *clear:
		if err := annotate.Clear(snap); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if err := saveSnapshot(*snapshotFile, snap); err != nil {
			fmt.Fprintf(os.Stderr, "error saving snapshot: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("annotation cleared")

	default:
		if *note == "" {
			fmt.Fprintln(os.Stderr, "error: --note is required when adding an annotation")
			os.Exit(1)
		}
		if err := annotate.Add(snap, *author, *note); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if err := saveSnapshot(*snapshotFile, snap); err != nil {
			fmt.Fprintf(os.Stderr, "error saving snapshot: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("annotation saved")
	}
}

// printAnnotation formats and prints the fields of an annotation to stdout.
func printAnnotation(a *annotate.Annotation) {
	fmt.Printf("Author : %s\n", a.Author)
	fmt.Printf("Note   : %s\n", a.Note)
	fmt.Printf("Date   : %s\n", a.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
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
	return os.WriteFile(path, data, 0o644)
}
