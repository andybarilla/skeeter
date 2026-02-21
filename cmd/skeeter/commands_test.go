package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (stdout string, stderr string, err error) {
	stdoutBuf := new(bytes.Buffer)
	stderrBuf := new(bytes.Buffer)

	root.SetOut(stdoutBuf)
	root.SetErr(stderrBuf)
	root.SetArgs(args)

	err = root.Execute()
	return stdoutBuf.String(), stderrBuf.String(), err
}

func setupTestEnv(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()
	repoDir := filepath.Join(dir, "repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatal(err)
	}

	oldWd, _ := os.Getwd()
	os.Chdir(repoDir)
	os.Unsetenv("SKEETER_DIR")

	cleanup := func() {
		os.Chdir(oldWd)
		os.Unsetenv("SKEETER_DIR")
	}

	return repoDir, cleanup
}

func TestInitCommand(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test-project")
	if err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	skeeterDir := ".skeeter"
	if _, err := os.Stat(skeeterDir); os.IsNotExist(err) {
		t.Error(".skeeter directory not created")
	}

	configPath := filepath.Join(skeeterDir, "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config.yaml not created")
	}

	tasksDir := filepath.Join(skeeterDir, "tasks")
	if _, err := os.Stat(tasksDir); os.IsNotExist(err) {
		t.Error("tasks directory not created")
	}

	templatesDir := filepath.Join(skeeterDir, "templates")
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		t.Error("templates directory not created")
	}
}

func TestCreateCommand(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "create", "Test task title", "-p", "high", "-t", "api,auth")
	if err != nil {
		t.Fatalf("create command failed: %v", err)
	}

	taskPath := filepath.Join(".skeeter", "tasks", "US-001.md")
	if _, err := os.Stat(taskPath); os.IsNotExist(err) {
		t.Fatal("task file not created")
	}

	content, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "title: Test task title") {
		t.Error("task file missing title")
	}
	if !strings.Contains(contentStr, "priority: high") {
		t.Error("task file missing priority")
	}
	if !strings.Contains(contentStr, "[api, auth]") {
		t.Error("task file missing tags")
	}
}

func TestCreateCommandWithInvalidPriority(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "create", "Test", "-p", "invalid-priority")
	if err == nil {
		t.Error("expected error for invalid priority")
	}
}

func TestListCommand(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "create", "Task 1", "-p", "high")
	if err != nil {
		t.Fatalf("create 1 failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "create", "Task 2", "-p", "low")
	if err != nil {
		t.Fatalf("create 2 failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "list")
	if err != nil {
		t.Fatalf("list command failed: %v", err)
	}

	if _, err := os.Stat(".skeeter/tasks/US-001.md"); os.IsNotExist(err) {
		t.Error("US-001.md not found")
	}
	if _, err := os.Stat(".skeeter/tasks/US-002.md"); os.IsNotExist(err) {
		t.Error("US-002.md not found")
	}
}

func TestListCommandEmpty(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "list")
	if err != nil {
		t.Fatalf("list command failed: %v", err)
	}
}

func TestListCommandWithFilter(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "create", "High Task", "-p", "high")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "create", "Low Task", "-p", "low")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "list", "--priority", "high")
	if err != nil {
		t.Fatalf("list command failed: %v", err)
	}
}

func TestShowCommand(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "create", "Test Task", "-p", "high")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "show", "US-001")
	if err != nil {
		t.Fatalf("show command failed: %v", err)
	}
}

func TestShowCommandNotFound(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "show", "US-999")
	if err == nil {
		t.Error("expected error for missing task")
	}
}

func TestStatusCommand(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "create", "Test Task")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "status", "US-001", "in-progress")
	if err != nil {
		t.Fatalf("status command failed: %v", err)
	}

	taskPath := filepath.Join(".skeeter", "tasks", "US-001.md")
	content, _ := os.ReadFile(taskPath)
	if !strings.Contains(string(content), "status: in-progress") {
		t.Error("task file not updated with new status")
	}
}

func TestStatusCommandInvalidStatus(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "create", "Test Task")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "status", "US-001", "invalid-status")
	if err == nil {
		t.Error("expected error for invalid status")
	}
}

func TestAssignCommand(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "create", "Test Task")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "assign", "US-001", "claude")
	if err != nil {
		t.Fatalf("assign command failed: %v", err)
	}

	taskPath := filepath.Join(".skeeter", "tasks", "US-001.md")
	content, _ := os.ReadFile(taskPath)
	if !strings.Contains(string(content), "assignee: claude") {
		t.Error("task file not updated with assignee")
	}
}

func TestNextCommandNoTasks(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "next")
	if err == nil {
		t.Error("expected error when no tasks available")
	}
}

func TestNextCommand(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "create", "Task 1", "-p", "high", "-s", "ready-for-development")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "create", "Task 2", "-p", "low", "-s", "ready-for-development")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "next")
	if err != nil {
		t.Fatalf("next command failed: %v", err)
	}
}

func TestNextCommandWithAssign(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "create", "Task 1", "-s", "ready-for-development")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "next", "--assign", "test-agent")
	if err != nil {
		t.Fatalf("next --assign failed: %v", err)
	}

	taskPath := filepath.Join(".skeeter", "tasks", "US-001.md")
	content, _ := os.ReadFile(taskPath)
	contentStr := string(content)

	if !strings.Contains(contentStr, "assignee: test-agent") {
		t.Error("task not assigned")
	}
	if !strings.Contains(contentStr, "status: in-progress") {
		t.Error("task status not changed to in-progress")
	}
}

func TestConfigCommand(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "config")
	if err != nil {
		t.Fatalf("config command failed: %v", err)
	}
}

func TestConfigSetCommand(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	_, _, err := executeCommand(rootCmd, "init", "test")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	_, _, err = executeCommand(rootCmd, "config", "set", "prefix", "TK")
	if err != nil {
		t.Fatalf("config set failed: %v", err)
	}

	configPath := filepath.Join(".skeeter", "config.yaml")
	content, _ := os.ReadFile(configPath)
	if !strings.Contains(string(content), "TK") {
		t.Error("config not updated with new prefix")
	}
}
