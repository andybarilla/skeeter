package task

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func Parse(content string) (*Task, error) {
	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "---") {
		return nil, fmt.Errorf("missing frontmatter delimiter")
	}

	// Remove leading ---
	rest := content[3:]
	frontmatter, body, found := strings.Cut(rest, "---")
	if !found {
		return nil, fmt.Errorf("missing closing frontmatter delimiter")
	}

	var t Task
	if err := yaml.Unmarshal([]byte(frontmatter), &t); err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	t.Body = strings.TrimRight(strings.TrimLeft(body, "\n"), "\n") + "\n"
	return &t, nil
}

func Marshal(t *Task) (string, error) {
	fm, err := yaml.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("marshaling frontmatter: %w", err)
	}

	var buf strings.Builder
	buf.WriteString("---\n")
	buf.Write(fm)
	buf.WriteString("---\n\n")
	buf.WriteString(t.Body)
	return buf.String(), nil
}
