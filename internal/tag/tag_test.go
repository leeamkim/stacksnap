package tag_test

import (
	"testing"

	"github.com/stacksnap/internal/snapshot"
	"github.com/stacksnap/internal/tag"
)

func makeSnap(tags ...string) *snapshot.Snapshot {
	return &snapshot.Snapshot{Tags: append([]string{}, tags...)}
}

func TestAdd_AppendsTags(t *testing.T) {
	s := makeSnap()
	if err := tag.Add(s, "backend", "go"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(s.Tags))
	}
}

func TestAdd_Deduplicates(t *testing.T) {
	s := makeSnap("go")
	_ = tag.Add(s, "go", "docker")
	if len(s.Tags) != 2 {
		t.Errorf("expected 2 tags after dedup, got %d", len(s.Tags))
	}
}

func TestAdd_InvalidTagReturnsError(t *testing.T) {
	s := makeSnap()
	if err := tag.Add(s, "bad tag!"); err == nil {
		t.Error("expected error for invalid tag name")
	}
}

func TestAdd_InvalidTagDoesNotMutate(t *testing.T) {
	s := makeSnap("go")
	_ = tag.Add(s, "bad tag!")
	if len(s.Tags) != 1 {
		t.Errorf("snapshot should not be mutated on error, got tags: %v", s.Tags)
	}
}

func TestRemove_DeletesTag(t *testing.T) {
	s := makeSnap("go", "docker", "backend")
	tag.Remove(s, "docker")
	for _, tg := range s.Tags {
		if tg == "docker" {
			t.Error("docker should have been removed")
		}
	}
	if len(s.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(s.Tags))
	}
}

func TestRemove_UnknownTagIgnored(t *testing.T) {
	s := makeSnap("go")
	tag.Remove(s, "nonexistent")
	if len(s.Tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(s.Tags))
	}
}

func TestList_ReturnsSorted(t *testing.T) {
	s := makeSnap("zebra", "apple", "mango")
	got := tag.List(s)
	want := []string{"apple", "mango", "zebra"}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("index %d: want %q got %q", i, v, got[i])
		}
	}
}

func TestList_EmptySnapshot(t *testing.T) {
	s := makeSnap()
	if got := tag.List(s); len(got) != 0 {
		t.Errorf("expected empty list, got %v", got)
	}
}
