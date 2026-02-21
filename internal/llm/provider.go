package llm

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// RunCLI pipes prompt to stdin of the configured command and captures stdout.
func RunCLI(ctx context.Context, command string, prompt string) (string, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", fmt.Errorf("no LLM command configured (run: skeeter config set llm.command \"claude -p\")")
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Stdin = strings.NewReader(prompt)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			return "", fmt.Errorf("llm command failed: %w\n%s", err, errMsg)
		}
		return "", fmt.Errorf("llm command failed: %w", err)
	}

	return strings.TrimSpace(stdout.String()), nil
}
