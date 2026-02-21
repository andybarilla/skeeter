package id

import (
	"os"
	"path/filepath"
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

func TestNextFromNamesDifferentPrefixes(t *testing.T) {
	names := []string{"US-001.md", "TK-001.md", "US-003.md", "TK-005.md"}
	got, err := NextFromNames(names, "US")
	if err != nil {
		t.Fatalf("NextFromNames: %v", err)
	}
	if got != "US-004" {
		t.Errorf("got %q, want %q", got, "US-004")
	}

	got, err = NextFromNames(names, "TK")
	if err != nil {
		t.Fatalf("NextFromNames: %v", err)
	}
	if got != "TK-006" {
		t.Errorf("got %q, want %q", got, "TK-006")
	}
}

func TestNextFromNamesInvalidNumbers(t *testing.T) {
	names := []string{"US-001.md", "US-abc.md", "US-.md", "US-003.md"}
	got, err := NextFromNames(names, "US")
	if err != nil {
		t.Fatalf("NextFromNames: %v", err)
	}
	if got != "US-004" {
		t.Errorf("got %q, want %q", got, "US-004")
	}
}

func TestNextFromNamesNoExtension(t *testing.T) {
	names := []string{"US-001", "US-002.md"}
	got, err := NextFromNames(names, "US")
	if err != nil {
		t.Fatalf("NextFromNames: %v", err)
	}
	if got != "US-003" {
		t.Errorf("got %q, want %q", got, "US-003")
	}
}

func TestNext(t *testing.T) {
	t.Run("nonexistent directory", func(t *testing.T) {
		got, err := Next("/nonexistent/path", "US")
		if err != nil {
			t.Fatalf("Next: %v", err)
		}
		if got != "US-001" {
			t.Errorf("got %q, want %q", got, "US-001")
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		dir := t.TempDir()
		got, err := Next(dir, "TK")
		if err != nil {
			t.Fatalf("Next: %v", err)
		}
		if got != "TK-001" {
			t.Errorf("got %q, want %q", got, "TK-001")
		}
	})

	t.Run("directory with tasks", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "US-001.md"), []byte(""), 0644)
		os.WriteFile(filepath.Join(dir, "US-002.md"), []byte(""), 0644)

		got, err := Next(dir, "US")
		if err != nil {
			t.Fatalf("Next: %v", err)
		}
		if got != "US-003" {
			t.Errorf("got %q, want %q", got, "US-003")
		}
	})

	t.Run("ignores subdirectories", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "US-001.md"), []byte(""), 0644)
		os.Mkdir(filepath.Join(dir, "US-002"), 0755)

		got, err := Next(dir, "US")
		if err != nil {
			t.Fatalf("Next: %v", err)
		}
		if got != "US-002" {
			t.Errorf("got %q, want %q", got, "US-002")
		}
	})
}
