// Package audit provides snapshot audit trail functionality,
// recording who accessed or modified a snapshot and when.
package audit

import (
	"errors"
	"fmt"
	"time"

	"github.com/stacksnap/internal/snapshot"
)

// EventKind describes the type of audit event.
type EventKind string

const (
	EventCapture  EventKind = "capture"
	EventExport   EventKind = "export"
	EventRestore  EventKind = "restore"
	EventValidate EventKind = "validate"
	EventShare    EventKind = "share"
)

// Event represents a single audit log entry.
type Event struct {
	Kind      EventKind `json:"kind"`
	Actor     string    `json:"actor"`
	Timestamp time.Time `json:"timestamp"`
	Note      string    `json:"note,omitempty"`
}

// Record appends an audit event to the snapshot's annotations.
// The actor may be empty, defaulting to "unknown".
func Record(snap *snapshot.Snapshot, kind EventKind, actor, note string) error {
	if snap == nil {
		return errors.New("audit: snapshot must not be nil")
	}
	if kind == "" {
		return errors.New("audit: event kind must not be empty")
	}
	if actor == "" {
		actor = "unknown"
	}
	if snap.Annotations == nil {
		snap.Annotations = make(map[string]string)
	}
	event := Event{
		Kind:      kind,
		Actor:     actor,
		Timestamp: time.Now().UTC(),
		Note:      note,
	}
	key := fmt.Sprintf("audit.%s.%s", kind, event.Timestamp.Format(time.RFC3339Nano))
	snap.Annotations[key] = fmt.Sprintf("actor=%s note=%s", actor, note)
	_ = event
	return nil
}

// List returns all audit events recorded on the snapshot.
func List(snap *snapshot.Snapshot) ([]Event, error) {
	if snap == nil {
		return nil, errors.New("audit: snapshot must not be nil")
	}
	var events []Event
	for k, v := range snap.Annotations {
		var kind EventKind
		var ts string
		if n, _ := fmt.Sscanf(k, "audit.%s", &kind); n != 1 {
			continue
		}
		// parse timestamp from key suffix
		if len(k) > len("audit.")+len(string(kind))+1 {
			ts = k[len("audit.")+len(string(kind))+1:]
		}
		t, _ := time.Parse(time.RFC3339Nano, ts)
		events = append(events, Event{
			Kind:      kind,
			Timestamp: t,
			Note:      v,
		})
	}
	return events, nil
}
