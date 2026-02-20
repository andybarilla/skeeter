package task

import (
	"strings"
	"testing"
)

func TestParseRoundtrip(t *testing.T) {
	input := `---
id: US-001
title: Test task
status: backlog
priority: high
assignee: andy
tags: [auth, api]
created: "2026-01-01"
updated: "2026-01-01"
---

This is the body.

## Acceptance Criteria

- [ ] Something
`

	tk, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if tk.ID != "US-001" {
		t.Errorf("ID = %q, want %q", tk.ID, "US-001")
	}
	if tk.Title != "Test task" {
		t.Errorf("Title = %q, want %q", tk.Title, "Test task")
	}
	if tk.Assignee != "andy" {
		t.Errorf("Assignee = %q, want %q", tk.Assignee, "andy")
	}
	if len(tk.Tags) != 2 {
		t.Errorf("Tags = %v, want 2 items", tk.Tags)
	}
	if !strings.Contains(tk.Body, "Acceptance Criteria") {
		t.Error("Body missing 'Acceptance Criteria'")
	}

	// Marshal and re-parse
	output, err := Marshal(tk)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	tk2, err := Parse(output)
	if err != nil {
		t.Fatalf("re-Parse: %v", err)
	}

	if tk2.ID != tk.ID || tk2.Title != tk.Title || tk2.Status != tk.Status {
		t.Errorf("roundtrip mismatch: got ID=%q Title=%q Status=%q", tk2.ID, tk2.Title, tk2.Status)
	}
}

func TestParseMissingFrontmatter(t *testing.T) {
	_, err := Parse("no frontmatter here")
	if err == nil {
		t.Error("expected error for missing frontmatter")
	}
}

func TestParseEmptyBody(t *testing.T) {
	input := `---
id: US-001
title: Empty body
status: backlog
created: "2026-01-01"
updated: "2026-01-01"
---
`
	tk, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if tk.ID != "US-001" {
		t.Errorf("ID = %q, want %q", tk.ID, "US-001")
	}
}
