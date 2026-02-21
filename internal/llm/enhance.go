package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/andybarilla/skeeter/internal/config"
	"github.com/andybarilla/skeeter/internal/task"
)

func EnhanceTask(ctx context.Context, cfg *config.Config, t *task.Task, template string) (string, error) {
	provider, err := NewProvider(cfg.LLM)
	if err != nil {
		return "", err
	}

	systemPrompt := buildSystemPrompt(cfg)
	userPrompt := buildUserPrompt(t.Title, t.Priority, t.Tags, t.Body, template)

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	return provider.Complete(ctx, messages)
}

func EnhanceDraft(ctx context.Context, cfg *config.Config, title, body, template string) (string, error) {
	provider, err := NewProvider(cfg.LLM)
	if err != nil {
		return "", err
	}

	systemPrompt := buildSystemPrompt(cfg)
	userPrompt := buildUserPrompt(title, "", nil, body, template)

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	return provider.Complete(ctx, messages)
}

func buildSystemPrompt(cfg *config.Config) string {
	var b strings.Builder
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
	b.WriteString("- Output ONLY the task body (markdown content), nothing else\n")

	return b.String()
}

func buildUserPrompt(title, priority string, tags []string, body, template string) string {
	var b strings.Builder

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
