// Package badge generates shield-style status badges summarising a snapshot.
package badge

import (
	"fmt"
	"strings"

	"github.com/stacksnap/internal/score"
	"github.com/stacksnap/internal/snapshot"
)

// Style controls the visual style of the generated badge.
type Style string

const (
	StyleFlat       Style = "flat"
	StyleFlatSquare Style = "flat-square"
	StylePlastic    Style = "plastic"
)

// Options configures badge generation.
type Options struct {
	Style  Style
	Label  string
	LogoID string // optional simple-icons slug
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Style: StyleFlat,
		Label: "stacksnap",
	}
}

// Result holds the generated badge artefacts.
type Result struct {
	ShieldsURL string // ready-to-embed shields.io URL
	Markdown   string // Markdown image snippet
	Grade      string
	Score      int
}

// Generate produces a shields.io badge URL for the given snapshot.
func Generate(snap *snapshot.Snapshot, opts Options) (*Result, error) {
	if snap == nil {
		return nil, fmt.Errorf("badge: snapshot must not be nil")
	}

	report, err := score.Evaluate(snap)
	if err != nil {
		return nil, fmt.Errorf("badge: scoring failed: %w", err)
	}

	color := colorFor(report.Grade)

	label := opts.Label
	if label == "" {
		label = "stacksnap"
	}
	style := opts.Style
	if style == "" {
		style = StyleFlat
	}

	message := fmt.Sprintf("%s%%20%d", report.Grade, report.Score)
	base := fmt.Sprintf("https://img.shields.io/badge/%s-%s-%s?style=%s",
		urlEncode(label), message, color, style)

	if opts.LogoID != "" {
		base += "&logo=" + opts.LogoID
	}

	md := fmt.Sprintf("![%s](%s)", label, base)

	return &Result{
		ShieldsURL: base,
		Markdown:   md,
		Grade:      report.Grade,
		Score:      report.Score,
	}, nil
}

func colorFor(grade string) string {
	switch strings.ToUpper(grade) {
	case "A":
		return "brightgreen"
	case "B":
		return "green"
	case "C":
		return "yellow"
	case "D":
		return "orange"
	default:
		return "red"
	}
}

func urlEncode(s string) string {
	return strings.ReplaceAll(s, " ", "%20")
}
