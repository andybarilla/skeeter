package llm

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/andybarilla/skeeter/internal/config"
	"github.com/andybarilla/skeeter/internal/task"
)

func TestBuildWorkPrompts(t *testing.T) {
	cfg := config.Default()
	cfg.Project.Name = "test-project"
	tk := &task.Task{
		ID:       "US-001",
		Title:    "Test Task",
		Priority: "high",
		Tags:     task.FlowSlice{"api", "auth"},
		Body:     "## Description\n\nDo the thing.\n",
	}

	t.Run("default prompts", func(t *testing.T) {
		tmpDir := t.TempDir()

		sys, user := BuildWorkPrompts(cfg, tk, tmpDir)

		if !strings.Contains(sys, "autonomous coding agent") {
			t.Error("system prompt missing expected content")
		}
		if !strings.Contains(sys, "test-project") {
			t.Error("system prompt missing project name")
		}
		if !strings.Contains(sys, tmpDir) {
			t.Error("system prompt missing skeeter dir")
		}

		if !strings.Contains(user, "US-001") {
			t.Error("user content missing task ID")
		}
		if !strings.Contains(user, "Test Task") {
			t.Error("user content missing task title")
		}
		if !strings.Contains(user, "high") {
			t.Error("user content missing priority")
		}
		if !strings.Contains(user, "api, auth") {
			t.Error("user content missing tags")
		}
		if !strings.Contains(user, "Do the thing") {
			t.Error("user content missing body")
		}
	})

	t.Run("custom prompt file", func(t *testing.T) {
		tmpDir := t.TempDir()
		promptsDir := filepath.Join(tmpDir, "prompts")
		if err := os.MkdirAll(promptsDir, 0755); err != nil {
			t.Fatal(err)
		}

		customPrompt := "Custom task: {{task_id}} - {{task_title}}"
		if err := os.WriteFile(filepath.Join(promptsDir, "work.md"), []byte(customPrompt), 0644); err != nil {
			t.Fatal(err)
		}

		_, user := BuildWorkPrompts(cfg, tk, tmpDir)

		if !strings.Contains(user, "Custom task: US-001 - Test Task") {
			t.Errorf("user content = %q, want custom prompt expanded", user)
		}
	})
}

func TestExpandPlaceholders(t *testing.T) {
	cfg := config.Default()
	cfg.Project.Name = "my-project"
	tk := &task.Task{
		ID:       "TK-123",
		Title:    "Placeholder Test",
		Priority: "critical",
		Tags:     task.FlowSlice{"tag1", "tag2"},
		Body:     "Body content here",
	}

	tests := []struct {
		template   string
		contains   string
		notContain string
	}{
		{"Task: {{task_id}}", "Task: TK-123", ""},
		{"Title: {{task_title}}", "Title: Placeholder Test", ""},
		{"Priority: {{task_priority}}", "Priority: critical", ""},
		{"Tags: {{task_tags}}", "Tags: tag1, tag2", ""},
		{"Body: {{task_body}}", "Body: Body content here", ""},
		{"Project: {{project_name}}", "Project: my-project", ""},
		{"Dir: {{skeeter_dir}}", "Dir: /some/path", ""},
	}

	for _, tt := range tests {
		result := expandPlaceholders(tt.template, cfg, tk, "/some/path")
		if tt.contains != "" && !strings.Contains(result, tt.contains) {
			t.Errorf("expandPlaceholders(%q) = %q, want it to contain %q", tt.template, result, tt.contains)
		}
		if tt.notContain != "" && strings.Contains(result, tt.notContain) {
			t.Errorf("expandPlaceholders(%q) = %q, want it NOT to contain %q", tt.template, result, tt.notContain)
		}
	}
}

func TestRunCLI(t *testing.T) {
	t.Run("successful execution reads stdin", func(t *testing.T) {
		tool := &config.LLMToolDef{
			Command:          "cat",
			PrintFlag:        "-",
			SystemPromptFlag: "",
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := RunCLI(ctx, tool, "system prompt", "user content")
		if err != nil {
			t.Fatalf("RunCLI: %v", err)
		}

		if !strings.Contains(result, "user content") {
			t.Errorf("result = %q, want it to contain user content", result)
		}
	})

	t.Run("command failure", func(t *testing.T) {
		tool := &config.LLMToolDef{
			Command:   "false",
			PrintFlag: "",
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := RunCLI(ctx, tool, "", "test")
		if err == nil {
			t.Error("expected error for failing command")
		}
	})

	t.Run("command not found", func(t *testing.T) {
		tool := &config.LLMToolDef{
			Command:   "nonexistent-command-xyz",
			PrintFlag: "",
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := RunCLI(ctx, tool, "", "test")
		if err == nil {
			t.Error("expected error for non-existent command")
		}
	})

	t.Run("empty command", func(t *testing.T) {
		tool := &config.LLMToolDef{
			Command:   "",
			PrintFlag: "",
		}

		ctx := context.Background()

		_, err := RunCLI(ctx, tool, "", "test")
		if err == nil {
			t.Error("expected error for empty command")
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		tool := &config.LLMToolDef{
			Command:   "sleep",
			PrintFlag: "",
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err := RunCLI(ctx, tool, "", "10")
		if err == nil {
			t.Error("expected error due to context cancellation")
		}
	})

	t.Run("stdin contains system prompt when no flag", func(t *testing.T) {
		tool := &config.LLMToolDef{
			Command:          "cat",
			PrintFlag:        "-",
			SystemPromptFlag: "",
		}

		ctx := context.Background()

		result, err := RunCLI(ctx, tool, "SYSTEM_PROMPT", "USER_CONTENT")
		if err != nil {
			t.Fatalf("RunCLI: %v", err)
		}

		if !strings.Contains(result, "SYSTEM_PROMPT") {
			t.Errorf("result = %q, want it to contain system prompt", result)
		}
		if !strings.Contains(result, "USER_CONTENT") {
			t.Errorf("result = %q, want it to contain user content", result)
		}
	})
}

func TestBuildArgs(t *testing.T) {
	tests := []struct {
		name        string
		tool        *config.LLMToolDef
		system      string
		extra       []string
		wantContain []string
	}{
		{
			name: "with system prompt flag",
			tool: &config.LLMToolDef{
				PrintFlag:        "-p",
				SystemPromptFlag: "--system",
			},
			system:      "sys",
			extra:       []string{"--extra"},
			wantContain: []string{"-p", "--system", "sys", "--extra"},
		},
		{
			name: "without system prompt flag",
			tool: &config.LLMToolDef{
				PrintFlag:        "-p",
				SystemPromptFlag: "",
			},
			system:      "sys",
			extra:       nil,
			wantContain: []string{"-p"},
		},
		{
			name: "empty system prompt",
			tool: &config.LLMToolDef{
				PrintFlag:        "-p",
				SystemPromptFlag: "--system",
			},
			system:      "",
			extra:       nil,
			wantContain: []string{"-p"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := buildArgs(tt.tool, tt.system, tt.extra)

			for _, want := range tt.wantContain {
				found := false
				for _, arg := range args {
					if arg == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("buildArgs() = %v, want it to contain %q", args, want)
				}
			}
		})
	}
}

func TestBuildStdin(t *testing.T) {
	tests := []struct {
		name     string
		tool     *config.LLMToolDef
		system   string
		user     string
		contains []string
	}{
		{
			name: "with system prompt flag",
			tool: &config.LLMToolDef{
				SystemPromptFlag: "--system",
			},
			system:   "SYSTEM",
			user:     "USER",
			contains: []string{"USER"},
		},
		{
			name: "without system prompt flag",
			tool: &config.LLMToolDef{
				SystemPromptFlag: "",
			},
			system:   "SYSTEM",
			user:     "USER",
			contains: []string{"SYSTEM", "USER"},
		},
		{
			name: "empty system prompt",
			tool: &config.LLMToolDef{
				SystemPromptFlag: "",
			},
			system:   "",
			user:     "USER",
			contains: []string{"USER"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildStdin(tt.tool, tt.system, tt.user)

			for _, want := range tt.contains {
				if !strings.Contains(result, want) {
					t.Errorf("buildStdin() = %q, want it to contain %q", result, want)
				}
			}
		})
	}
}

func TestCleanLLMEnv(t *testing.T) {
	originalClaudeCode := os.Getenv("CLAUDECODE")
	defer os.Setenv("CLAUDECODE", originalClaudeCode)

	os.Setenv("CLAUDECODE", "test-value")
	os.Setenv("OTHER_VAR", "should-remain")

	env := cleanLLMEnv()

	for _, e := range env {
		if strings.HasPrefix(e, "CLAUDECODE=") {
			t.Error("cleanLLMEnv should remove CLAUDECODE")
		}
	}

	found := false
	for _, e := range env {
		if strings.HasPrefix(e, "OTHER_VAR=") {
			found = true
			break
		}
	}
	if !found {
		t.Error("cleanLLMEnv should keep OTHER_VAR")
	}
}

func TestEnhanceSystemPrompt(t *testing.T) {
	cfg := config.Default()
	cfg.Project.Name = "enhance-test"

	result := enhanceSystemPrompt(cfg)

	if !strings.Contains(result, "technical project manager") {
		t.Error("missing role description")
	}
	if !strings.Contains(result, "enhance-test") {
		t.Error("missing project name")
	}
	if !strings.Contains(result, "backlog") {
		t.Error("missing statuses")
	}
	if !strings.Contains(result, "critical") {
		t.Error("missing priorities")
	}
}

func TestEnhanceUserContent(t *testing.T) {
	t.Run("with all fields", func(t *testing.T) {
		result := enhanceUserContent(
			"Task Title",
			"high",
			[]string{"tag1", "tag2"},
			"Existing body content",
			"Template content",
		)

		if !strings.Contains(result, "Task Title") {
			t.Error("missing title")
		}
		if !strings.Contains(result, "high") {
			t.Error("missing priority")
		}
		if !strings.Contains(result, "tag1, tag2") {
			t.Error("missing tags")
		}
		if !strings.Contains(result, "Template content") {
			t.Error("missing template")
		}
		if !strings.Contains(result, "Existing body content") {
			t.Error("missing existing body")
		}
	})

	t.Run("minimal fields", func(t *testing.T) {
		result := enhanceUserContent(
			"Simple Task",
			"",
			nil,
			"",
			"Template",
		)

		if !strings.Contains(result, "Simple Task") {
			t.Error("missing title")
		}
		if strings.Contains(result, "Existing notes") {
			t.Error("should not include existing notes section when body is empty")
		}
	})

	t.Run("body same as template not duplicated", func(t *testing.T) {
		template := "## Acceptance Criteria\n\n- [ ]\n"
		result := enhanceUserContent(
			"Task",
			"",
			nil,
			template,
			template,
		)

		if strings.Contains(result, "Existing notes") {
			t.Error("should not include existing notes when body equals template")
		}
	})
}

func TestDefaultWorkUserContent(t *testing.T) {
	cfg := config.Default()
	tk := &task.Task{
		ID:       "US-099",
		Title:    "Work Task",
		Priority: "medium",
		Tags:     task.FlowSlice{"frontend"},
		Body:     "## Description\n\nBuild the UI.\n",
	}

	result := defaultWorkUserContent(cfg, tk, "/path/to/.skeeter")

	if !strings.Contains(result, "US-099") {
		t.Error("missing task ID")
	}
	if !strings.Contains(result, "Work Task") {
		t.Error("missing task title")
	}
	if !strings.Contains(result, "medium") {
		t.Error("missing priority")
	}
	if !strings.Contains(result, "frontend") {
		t.Error("missing tags")
	}
	if !strings.Contains(result, "Build the UI") {
		t.Error("missing body")
	}
	if !strings.Contains(result, "### Instructions") {
		t.Error("missing instructions section")
	}
}
