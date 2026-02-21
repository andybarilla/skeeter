package llm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andybarilla/skeeter/internal/config"
	"github.com/andybarilla/skeeter/internal/task"
)

// BuildWorkPrompt constructs the prompt sent to the LLM work command.
// If .skeeter/prompts/work.md exists, it is used as the template with placeholder substitution.
// Otherwise a default prompt is generated.
func BuildWorkPrompt(cfg *config.Config, t *task.Task, skeeterDir string) string {
	customPath := filepath.Join(skeeterDir, "prompts", "work.md")
	if data, err := os.ReadFile(customPath); err == nil {
		return expandPlaceholders(string(data), cfg, t, skeeterDir)
	}
	return defaultWorkPrompt(cfg, t, skeeterDir)
}

func expandPlaceholders(tmpl string, cfg *config.Config, t *task.Task, skeeterDir string) string {
	r := strings.NewReplacer(
		"{{task_id}}", t.ID,
		"{{task_title}}", t.Title,
		"{{task_priority}}", t.Priority,
		"{{task_tags}}", strings.Join(t.Tags, ", "),
		"{{task_body}}", t.Body,
		"{{project_name}}", cfg.Project.Name,
		"{{skeeter_dir}}", skeeterDir,
	)
	return r.Replace(tmpl)
}

func defaultWorkPrompt(cfg *config.Config, t *task.Task, skeeterDir string) string {
	var b strings.Builder

	b.WriteString("You are an autonomous coding agent. Implement the following task.\n\n")

	if cfg.Project.Name != "" {
		fmt.Fprintf(&b, "Project: %s\n", cfg.Project.Name)
	}
	fmt.Fprintf(&b, "Skeeter directory: %s\n\n", skeeterDir)

	fmt.Fprintf(&b, "## Task %s: %s\n\n", t.ID, t.Title)

	if t.Priority != "" {
		fmt.Fprintf(&b, "Priority: %s\n", t.Priority)
	}
	if len(t.Tags) > 0 {
		fmt.Fprintf(&b, "Tags: %s\n", strings.Join(t.Tags, ", "))
	}

	if t.Body != "" {
		b.WriteString("\n### Description\n\n")
		b.WriteString(t.Body)
		b.WriteString("\n")
	}

	b.WriteString("\n### Instructions\n\n")
	b.WriteString("- Implement the task described above following existing code patterns and conventions.\n")
	b.WriteString("- Run tests and/or build to validate your changes compile and pass.\n")
	b.WriteString("- Commit your changes with a message referencing the task ID (e.g., \"" + t.ID + ": <summary>\").\n")
	b.WriteString("- Do NOT modify any files inside the " + skeeterDir + "/tasks/ directory.\n")

	return b.String()
}
