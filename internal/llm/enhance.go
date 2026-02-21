package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/andybarilla/skeeter/internal/config"
	"github.com/andybarilla/skeeter/internal/task"
)

func EnhanceTask(ctx context.Context, cfg *config.Config, t *task.Task, template string) (string, error) {
	if cfg.LLM.Command == "" {
		return "", fmt.Errorf("no LLM command configured (run: skeeter config set llm.command \"claude -p\")")
	}

	prompt := buildPrompt(cfg, t.Title, t.Priority, t.Tags, t.Body, template)
	return RunCLI(ctx, cfg.LLM.Command, prompt)
}

func EnhanceDraft(ctx context.Context, cfg *config.Config, title, body, template string) (string, error) {
	if cfg.LLM.Command == "" {
		return "", fmt.Errorf("no LLM command configured (run: skeeter config set llm.command \"claude -p\")")
	}

	prompt := buildPrompt(cfg, title, "", nil, body, template)
	return RunCLI(ctx, cfg.LLM.Command, prompt)
}

func buildPrompt(cfg *config.Config, title, priority string, tags []string, body, template string) string {
	var b strings.Builder

	// System context
	b.WriteString("You are a technical project manager helping flesh out task descriptions for a software project.\n\n")

	if cfg.Project.Name != "" {
		fmt.Fprintf(&b, "Project: %s\n", cfg.Project.Name)
	}
	fmt.Fprintf(&b, "Workflow statuses: %s\n", strings.Join(cfg.Statuses, " -> "))
	fmt.Fprintf(&b, "Priority levels: %s\n\n", strings.Join(cfg.Priorities, ", "))

	b.WriteString("Guidelines:\n")
	b.WriteString("- Follow the template structure provided, filling in each section\n")
	b.WriteString("- Write concrete, testable acceptance criteria as checklist items\n")
	b.WriteString("- Add relevant technical context and considerations\n")
	b.WriteString("- Keep the language clear and actionable\n")
	b.WriteString("- Do NOT include YAML frontmatter in your output\n")
	b.WriteString("- Output ONLY the task body (markdown content), nothing else\n\n")

	// Task details
	fmt.Fprintf(&b, "Task title: %s\n", title)
	if priority != "" {
		fmt.Fprintf(&b, "Priority: %s\n", priority)
	}
	if len(tags) > 0 {
		fmt.Fprintf(&b, "Tags: %s\n", strings.Join(tags, ", "))
	}

	b.WriteString("\nTemplate structure to follow:\n```\n")
	b.WriteString(template)
	b.WriteString("```\n")

	trimmedBody := strings.TrimSpace(body)
	trimmedTemplate := strings.TrimSpace(template)
	if trimmedBody != "" && trimmedBody != trimmedTemplate {
		b.WriteString("\nExisting notes/body:\n```\n")
		b.WriteString(body)
		b.WriteString("```\n")
	}

	b.WriteString("\nPlease flesh out this task following the template structure. Write detailed acceptance criteria and add relevant context.")

	return b.String()
}
