// Package watch monitors the local dev environment for changes relative to
// a previously captured snapshot, emitting a diff whenever a tool version or
// relevant environment variable changes.
package watch

import (
	"context"
	"fmt"
	"time"

	"github.com/stacksnap/stacksnap/internal/diff"
	"github.com/stacksnap/stacksnap/internal/snapshot"
)

// DefaultInterval is the polling interval used when none is specified.
const DefaultInterval = 30 * time.Second

// Options controls the behaviour of Watch.
type Options struct {
	// Interval between successive environment polls.
	Interval time.Duration

	// OnChange is called every time a difference is detected relative to base.
	// The callback receives a human-readable summary of what changed.
	// If OnChange is nil a default printer to stdout is used.
	OnChange func(summary string)

	// OnError is called when a poll cycle fails. If nil, errors are silently
	// discarded so that transient failures do not stop the watcher.
	OnError func(err error)
}

// DefaultOptions returns an Options value with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Interval: DefaultInterval,
		OnChange: func(summary string) {
			fmt.Println("[stacksnap watch]", summary)
		},
		OnError: func(err error) {
			fmt.Println("[stacksnap watch] poll error:", err)
		},
	}
}

// Watch polls the local environment at the configured interval and invokes
// opts.OnChange whenever the live state diverges from base. It blocks until
// ctx is cancelled, at which point it returns ctx.Err().
func Watch(ctx context.Context, base *snapshot.Snapshot, opts Options) error {
	if base == nil {
		return fmt.Errorf("watch: base snapshot must not be nil")
	}
	if opts.Interval <= 0 {
		opts.Interval = DefaultInterval
	}
	if opts.OnChange == nil {
		opts.OnChange = DefaultOptions().OnChange
	}
	if opts.OnError == nil {
		opts.OnError = func(error) {} // no-op
	}

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := poll(base, opts); err != nil {
				opts.OnError(err)
			}
		}
	}
}

// poll captures a fresh snapshot and compares it against base. If any
// differences are found it invokes opts.OnChange with the summary text.
func poll(base *snapshot.Snapshot, opts Options) error {
	live, err := snapshot.Capture()
	if err != nil {
		return fmt.Errorf("watch: capture failed: %w", err)
	}

	result := diff.Compare(base, live)
	if len(result.Added)+len(result.Removed)+len(result.Changed) == 0 {
		// No change — nothing to report.
		return nil
	}

	opts.OnChange(diff.Summary(result))
	return nil
}
