package llm

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/andybarilla/skeeter/internal/config"
)

// cleanLLMEnv returns a copy of the current environment with variables removed
// that prevent LLM CLIs from running as subprocesses (e.g. CLAUDECODE prevents
// nested Claude Code sessions).
func cleanLLMEnv() []string {
	var env []string
	for _, e := range os.Environ() {
		key, _, _ := strings.Cut(e, "=")
		if key == "CLAUDECODE" {
			continue
		}
		env = append(env, e)
	}
	return env
}

// buildArgs constructs the argument list for an LLM CLI invocation.
func buildArgs(tool *config.LLMToolDef, systemPrompt string, extraArgs []string) []string {
	var args []string
	args = append(args, tool.PrintFlag)
	if tool.SystemPromptFlag != "" && systemPrompt != "" {
		args = append(args, tool.SystemPromptFlag, systemPrompt)
	}
	args = append(args, extraArgs...)
	return args
}

// buildStdin returns the content to pipe to stdin. If the tool has no
// SystemPromptFlag, the system prompt is prepended to the user content.
func buildStdin(tool *config.LLMToolDef, systemPrompt, userContent string) string {
	if tool.SystemPromptFlag != "" || systemPrompt == "" {
		return userContent
	}
	return systemPrompt + "\n\n" + userContent
}

// RunCLI invokes the tool with separated system prompt and user content, capturing stdout.
func RunCLI(ctx context.Context, tool *config.LLMToolDef, systemPrompt, userContent string, extraArgs ...string) (string, error) {
	if tool.Command == "" {
		return "", fmt.Errorf("no LLM tool command configured (run: skeeter config set llm.tool claude)")
	}

	args := buildArgs(tool, systemPrompt, extraArgs)
	cmd := exec.CommandContext(ctx, tool.Command, args...)
	cmd.Stdin = strings.NewReader(buildStdin(tool, systemPrompt, userContent))
	cmd.Env = cleanLLMEnv()

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

// RunCLIPassthrough invokes the tool with separated system prompt and user content,
// wiring stdout/stderr to the terminal.
func RunCLIPassthrough(ctx context.Context, tool *config.LLMToolDef, systemPrompt, userContent string, extraArgs ...string) error {
	if tool.Command == "" {
		return fmt.Errorf("no LLM tool command configured (run: skeeter config set llm.tool claude)")
	}

	args := buildArgs(tool, systemPrompt, extraArgs)
	cmd := exec.CommandContext(ctx, tool.Command, args...)
	cmd.Stdin = strings.NewReader(buildStdin(tool, systemPrompt, userContent))
	cmd.Env = cleanLLMEnv()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
