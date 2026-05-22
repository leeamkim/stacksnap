// Command share encodes an existing snapshot file into a shareable string
// or decodes a share string back to a snapshot JSON file.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/user/stacksnap/internal/share"
	"github.com/user/stacksnap/internal/snapshot"
)

func main() {
	encodeCmd := flag.NewFlagSet("encode", flag.ExitOnError)
	encodeInput := encodeCmd.String("snapshot", "", "path to snapshot JSON file (required)")

	decodeCmd := flag.NewFlagSet("decode", flag.ExitOnError)
	decodeStr := decodeCmd.String("share", "", "share string to decode (required)")
	decodeOut := decodeCmd.String("out", "decoded-snapshot.json", "output file path")

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: share <encode|decode> [flags]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "encode":
		_ = encodeCmd.Parse(os.Args[2:])
		if *encodeInput == "" {
			fmt.Fprintln(os.Stderr, "error: -snapshot is required")
			os.Exit(1)
		}
		snap, err := loadSnapshot(*encodeInput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading snapshot: %v\n", err)
			os.Exit(1)
		}
		result, err := share.Encode(snap)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error encoding: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(result)

	case "decode":
		_ = decodeCmd.Parse(os.Args[2:])
		if *decodeStr == "" {
			fmt.Fprintln(os.Stderr, "error: -share is required")
			os.Exit(1)
		}
		snap, err := share.Decode(*decodeStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error decoding: %v\n", err)
			os.Exit(1)
		}
		data, _ := json.MarshalIndent(snap, "", "  ")
		if err := os.WriteFile(*decodeOut, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("snapshot written to %s\n", *decodeOut)

	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func loadSnapshot(path string) (snapshot.Snapshot, error) {
	var snap snapshot.Snapshot
	data, err := os.ReadFile(path)
	if err != nil {
		return snap, err
	}
	err = json.Unmarshal(data, &snap)
	return snap, err
}
