package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Project.Name != "" {
		t.Errorf("Project.Name = %q, want empty", cfg.Project.Name)
	}
	if cfg.Project.Prefix != "US" {
		t.Errorf("Project.Prefix = %q, want %q", cfg.Project.Prefix, "US")
	}

	expectedStatuses := []string{"backlog", "ready-for-development", "in-progress", "done"}
	if len(cfg.Statuses) != len(expectedStatuses) {
		t.Errorf("Statuses length = %d, want %d", len(cfg.Statuses), len(expectedStatuses))
	}
	for i, s := range expectedStatuses {
		if cfg.Statuses[i] != s {
			t.Errorf("Statuses[%d] = %q, want %q", i, cfg.Statuses[i], s)
		}
	}

	expectedPriorities := []string{"critical", "high", "medium", "low"}
	if len(cfg.Priorities) != len(expectedPriorities) {
		t.Errorf("Priorities length = %d, want %d", len(cfg.Priorities), len(expectedPriorities))
	}
	for i, p := range expectedPriorities {
		if cfg.Priorities[i] != p {
			t.Errorf("Priorities[%d] = %q, want %q", i, cfg.Priorities[i], p)
		}
	}

	if cfg.AutoCommit {
		t.Error("AutoCommit = true, want false")
	}

	if cfg.LLM.Tool != "claude" {
		t.Errorf("LLM.Tool = %q, want %q", cfg.LLM.Tool, "claude")
	}
}

func TestLoad(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		dir := t.TempDir()
		configContent := `project:
  name: test-project
  prefix: TK
statuses:
  - todo
  - doing
  - done
priorities:
  - p0
  - p1
  - p2
auto_commit: true
llm:
  tool: custom
`
		err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("writing config: %v", err)
		}

		cfg, err := Load(dir)
		if err != nil {
			t.Fatalf("Load: %v", err)
		}

		if cfg.Project.Name != "test-project" {
			t.Errorf("Project.Name = %q, want %q", cfg.Project.Name, "test-project")
		}
		if cfg.Project.Prefix != "TK" {
			t.Errorf("Project.Prefix = %q, want %q", cfg.Project.Prefix, "TK")
		}
		if len(cfg.Statuses) != 3 || cfg.Statuses[0] != "todo" {
			t.Errorf("Statuses = %v, want [todo doing done]", cfg.Statuses)
		}
		if len(cfg.Priorities) != 3 || cfg.Priorities[0] != "p0" {
			t.Errorf("Priorities = %v, want [p0 p1 p2]", cfg.Priorities)
		}
		if !cfg.AutoCommit {
			t.Error("AutoCommit = false, want true")
		}
		if cfg.LLM.Tool != "custom" {
			t.Errorf("LLM.Tool = %q, want %q", cfg.LLM.Tool, "custom")
		}
	})

	t.Run("partial config uses defaults", func(t *testing.T) {
		dir := t.TempDir()
		configContent := `project:
  name: partial
`
		err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("writing config: %v", err)
		}

		cfg, err := Load(dir)
		if err != nil {
			t.Fatalf("Load: %v", err)
		}

		if cfg.Project.Name != "partial" {
			t.Errorf("Project.Name = %q, want %q", cfg.Project.Name, "partial")
		}
		if cfg.Project.Prefix != "US" {
			t.Errorf("Project.Prefix = %q, want default %q", cfg.Project.Prefix, "US")
		}
		if len(cfg.Statuses) != 4 {
			t.Errorf("Statuses = %v, want default 4 statuses", cfg.Statuses)
		}
	})

	t.Run("missing config file", func(t *testing.T) {
		dir := t.TempDir()
		_, err := Load(dir)
		if err == nil {
			t.Error("expected error for missing config file")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		dir := t.TempDir()
		err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("invalid: [yaml"), 0644)
		if err != nil {
			t.Fatalf("writing config: %v", err)
		}

		_, err = Load(dir)
		if err == nil {
			t.Error("expected error for invalid YAML")
		}
	})
}

func TestSave(t *testing.T) {
	t.Run("save and reload", func(t *testing.T) {
		dir := t.TempDir()
		cfg := Default()
		cfg.Project.Name = "save-test"
		cfg.Project.Prefix = "SV"
		cfg.AutoCommit = true

		err := cfg.Save(dir)
		if err != nil {
			t.Fatalf("Save: %v", err)
		}

		loaded, err := Load(dir)
		if err != nil {
			t.Fatalf("Load after Save: %v", err)
		}

		if loaded.Project.Name != "save-test" {
			t.Errorf("Project.Name = %q, want %q", loaded.Project.Name, "save-test")
		}
		if loaded.Project.Prefix != "SV" {
			t.Errorf("Project.Prefix = %q, want %q", loaded.Project.Prefix, "SV")
		}
		if !loaded.AutoCommit {
			t.Error("AutoCommit = false, want true")
		}
	})
}

func TestValidStatus(t *testing.T) {
	cfg := Default()

	tests := []struct {
		status   string
		expected bool
	}{
		{"backlog", true},
		{"ready-for-development", true},
		{"in-progress", true},
		{"done", true},
		{"unknown", false},
		{"", false},
		{"BACKLOG", false},
	}

	for _, tt := range tests {
		result := cfg.ValidStatus(tt.status)
		if result != tt.expected {
			t.Errorf("ValidStatus(%q) = %v, want %v", tt.status, result, tt.expected)
		}
	}
}

func TestValidPriority(t *testing.T) {
	cfg := Default()

	tests := []struct {
		priority string
		expected bool
	}{
		{"critical", true},
		{"high", true},
		{"medium", true},
		{"low", true},
		{"unknown", false},
		{"", false},
		{"CRITICAL", false},
	}

	for _, tt := range tests {
		result := cfg.ValidPriority(tt.priority)
		if result != tt.expected {
			t.Errorf("ValidPriority(%q) = %v, want %v", tt.priority, result, tt.expected)
		}
	}
}

func TestPriorityRank(t *testing.T) {
	cfg := Default()

	tests := []struct {
		priority string
		expected int
	}{
		{"critical", 0},
		{"high", 1},
		{"medium", 2},
		{"low", 3},
		{"unknown", 4},
		{"", 4},
	}

	for _, tt := range tests {
		result := cfg.PriorityRank(tt.priority)
		if result != tt.expected {
			t.Errorf("PriorityRank(%q) = %d, want %d", tt.priority, result, tt.expected)
		}
	}
}

func TestResolveTool(t *testing.T) {
	t.Run("builtin claude", func(t *testing.T) {
		cfg := Default()
		cfg.LLM.Tool = "claude"

		tool, err := cfg.ResolveTool()
		if err != nil {
			t.Fatalf("ResolveTool: %v", err)
		}
		if tool.Command != "claude" {
			t.Errorf("Command = %q, want %q", tool.Command, "claude")
		}
		if tool.PrintFlag != "-p" {
			t.Errorf("PrintFlag = %q, want %q", tool.PrintFlag, "-p")
		}
	})

	t.Run("custom tool", func(t *testing.T) {
		cfg := Default()
		cfg.LLM.Tool = "my-tool"
		cfg.LLM.Tools = map[string]LLMToolDef{
			"my-tool": {
				Command:   "my-llm",
				PrintFlag: "--print",
			},
		}

		tool, err := cfg.ResolveTool()
		if err != nil {
			t.Fatalf("ResolveTool: %v", err)
		}
		if tool.Command != "my-llm" {
			t.Errorf("Command = %q, want %q", tool.Command, "my-llm")
		}
		if tool.PrintFlag != "--print" {
			t.Errorf("PrintFlag = %q, want %q", tool.PrintFlag, "--print")
		}
	})

	t.Run("no tool configured", func(t *testing.T) {
		cfg := Default()
		cfg.LLM.Tool = ""

		_, err := cfg.ResolveTool()
		if err == nil {
			t.Error("expected error for empty tool")
		}
	})

	t.Run("unknown tool", func(t *testing.T) {
		cfg := Default()
		cfg.LLM.Tool = "nonexistent"

		_, err := cfg.ResolveTool()
		if err == nil {
			t.Error("expected error for unknown tool")
		}
	})
}
