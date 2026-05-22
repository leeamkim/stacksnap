package redact

import (
	"regexp"
	"strings"

	"github.com/stacksnap/internal/snapshot"
)

// DefaultPatterns are common sensitive key patterns to redact.
var DefaultPatterns = []string{
	"(?i)password",
	"(?i)secret",
	"(?i)token",
	"(?i)api[_-]?key",
	"(?i)private[_-]?key",
	"(?i)auth",
	"(?i)credential",
}

const redactedValue = "[REDACTED]"

// Options controls redaction behaviour.
type Options struct {
	// Patterns is a list of regex patterns matched against env var keys.
	Patterns []string
	// Placeholder overrides the default redacted placeholder string.
	Placeholder string
}

// Apply returns a shallow copy of snap with matching env values redacted.
// The original snapshot is never mutated.
func Apply(snap *snapshot.Snapshot, opts *Options) (*snapshot.Snapshot, error) {
	if snap == nil {
		return nil, fmt.Errorf("redact: snapshot must not be nil")
	}

	if opts == nil {
		opts = &Options{Patterns: DefaultPatterns}
	}
	if len(opts.Patterns) == 0 {
		opts.Patterns = DefaultPatterns
	}
	placeholder := redactedValue
	if opts.Placeholder != "" {
		placeholder = opts.Placeholder
	}

	compiled := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("redact: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}

	newEnv := make(map[string]string, len(snap.Env))
	for k, v := range snap.Env {
		if matchesAny(k, compiled) {
			newEnv[k] = placeholder
		} else {
			newEnv[k] = v
		}
	}

	copy := *snap
	copy.Env = newEnv
	return &copy, nil
}

func matchesAny(key string, patterns []*regexp.Regexp) bool {
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// SensitiveKeys returns the list of env keys that would be redacted.
func SensitiveKeys(snap *snapshot.Snapshot, opts *Options) []string {
	if snap == nil || len(snap.Env) == 0 {
		return nil
	}
	if opts == nil {
		opts = &Options{Patterns: DefaultPatterns}
	}
	compiled := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			continue
		}
		compiled = append(compiled, re)
	}
	var keys []string
	for k := range snap.Env {
		if matchesAny(k, compiled) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}
