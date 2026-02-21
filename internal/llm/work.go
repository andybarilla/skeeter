package llm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andybarilla/skeeter/internal/config"
	"github.com/andybarilla/skeeter/internal/task"
)

// BuildWorkPrompts constructs separated system prompt and user content for the work command.
// If .skeeter/prompts/work.md exists, it is used as user content with placeholder substitution.
// Otherwise default prompts are generated.
func BuildWorkPrompts(cfg *config.Config, t *task.Task, skeeterDir string) (systemPrompt, userContent string) {
	systemPrompt = workSystemPrompt(cfg, skeeterDir)

	customPath := filepath.Join(skeeterDir, "prompts", "work.md")
	if data, err := os.ReadFile(customPath); err == nil {
		userContent = expandPlaceholders(string(data), cfg, t, skeeterDir)
	} else {
		userContent = defaultWorkUserContent(cfg, t, skeeterDir)
	}
	return
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

func workSystemPrompt(cfg *config.Config, skeeterDir string) string {
	var b strings.Builder

	b.WriteString("You are an autonomous coding agent. Implement the task provided by the user.")

	if cfg.Project.Name != "" {
		fmt.Fprintf(&b, "\n\nProject: %s", cfg.Project.Name)
	}
	fmt.Fprintf(&b, "\nSkeeter directory: %s", skeeterDir)

	return b.String()
}

func defaultWorkUserContent(cfg *config.Config, t *task.Task, skeeterDir string) string {
	var b strings.Builder

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
