package id

import (
	"testing"
)

func TestNextFromNamesEmpty(t *testing.T) {
	got, err := NextFromNames(nil, "US")
	if err != nil {
		t.Fatalf("NextFromNames: %v", err)
	}
	if got != "US-001" {
		t.Errorf("got %q, want %q", got, "US-001")
	}
}

func TestNextFromNamesSequential(t *testing.T) {
	names := []string{"US-001.md", "US-002.md", "US-003.md"}
	got, err := NextFromNames(names, "US")
	if err != nil {
		t.Fatalf("NextFromNames: %v", err)
	}
	if got != "US-004" {
		t.Errorf("got %q, want %q", got, "US-004")
	}
}

func TestNextFromNamesGap(t *testing.T) {
	names := []string{"US-001.md", "US-005.md"}
	got, err := NextFromNames(names, "US")
	if err != nil {
		t.Fatalf("NextFromNames: %v", err)
	}
	if got != "US-006" {
		t.Errorf("got %q, want %q (should use max+1, not fill gaps)", got, "US-006")
	}
}

func TestNextFromNamesCustomPrefix(t *testing.T) {
	names := []string{"TASK-001.md", "TASK-002.md"}
	got, err := NextFromNames(names, "TASK")
	if err != nil {
		t.Fatalf("NextFromNames: %v", err)
	}
	if got != "TASK-003" {
		t.Errorf("got %q, want %q", got, "TASK-003")
	}
}

func TestNextFromNamesIgnoresNonMatching(t *testing.T) {
	names := []string{"US-001.md", "README.md", "notes.txt", "US-003.md"}
	got, err := NextFromNames(names, "US")
	if err != nil {
		t.Fatalf("NextFromNames: %v", err)
	}
	if got != "US-004" {
		t.Errorf("got %q, want %q", got, "US-004")
	}
}
