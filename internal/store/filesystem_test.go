package store

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andybarilla/skeeter/internal/task"
)

func setupTestStore(t *testing.T) *FilesystemStore {
	t.Helper()
	dir := t.TempDir()
	skeeterDir := filepath.Join(dir, ".skeeter")

	s := &FilesystemStore{Dir: skeeterDir}
	if err := s.Init("test-project"); err != nil {
		t.Fatalf("Init: %v", err)
	}
	return s
}

func setupTestStoreWithGit(t *testing.T) (*FilesystemStore, string) {
	t.Helper()
	dir := t.TempDir()

	// Initialize git repo
	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}

	skeeterDir := filepath.Join(dir, ".skeeter")
	s := &FilesystemStore{Dir: skeeterDir}
	if err := s.Init("test-project"); err != nil {
		t.Fatalf("Init: %v", err)
	}

	// Initial commit so auto-commit has something to build on
	cmd := exec.Command("git", "add", "-A")
	cmd.Dir = dir
	cmd.Run()
	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = dir
	cmd.Run()

	return s, dir
}

func gitLog(t *testing.T, dir string) []string {
	t.Helper()
	cmd := exec.Command("git", "log", "--oneline", "--format=%s")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git log: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return lines
}

func TestInit(t *testing.T) {
	s := setupTestStore(t)

	// Verify directories exist
	if _, err := os.Stat(filepath.Join(s.Dir, "tasks")); err != nil {
		t.Errorf("tasks dir not created: %v", err)
	}
	if _, err := os.Stat(filepath.Join(s.Dir, "templates")); err != nil {
		t.Errorf("templates dir not created: %v", err)
	}

	// Verify config
	if s.Config.Project.Name != "test-project" {
		t.Errorf("project name = %q, want %q", s.Config.Project.Name, "test-project")
	}
	if s.Config.Project.Prefix != "US" {
		t.Errorf("prefix = %q, want %q", s.Config.Project.Prefix, "US")
	}

	// Verify SKEETER.md exists
	if _, err := os.Stat(filepath.Join(s.Dir, "SKEETER.md")); err != nil {
		t.Errorf("SKEETER.md not created: %v", err)
	}

	// Verify default template exists
	if _, err := os.Stat(filepath.Join(s.Dir, "templates", "default.md")); err != nil {
		t.Errorf("default template not created: %v", err)
	}
}

func TestCreateAndGet(t *testing.T) {
	s := setupTestStore(t)

	tk := &task.Task{
		ID:       "US-001",
		Title:    "Test task",
		Status:   "backlog",
		Priority: "high",
		Tags:     task.FlowSlice{"auth", "api"},
		Created:  "2026-01-01",
		Updated:  "2026-01-01",
		Body:     "Task body here.\n",
	}

	if err := s.Create(tk); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := s.Get("US-001")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if got.Title != "Test task" {
		t.Errorf("Title = %q, want %q", got.Title, "Test task")
	}
	if got.Status != "backlog" {
		t.Errorf("Status = %q, want %q", got.Status, "backlog")
	}
	if got.Priority != "high" {
		t.Errorf("Priority = %q, want %q", got.Priority, "high")
	}
	if len(got.Tags) != 2 || got.Tags[0] != "auth" {
		t.Errorf("Tags = %v, want [auth, api]", got.Tags)
	}
	if got.Body != "Task body here.\n" {
		t.Errorf("Body = %q, want %q", got.Body, "Task body here.\n")
	}
}

func TestGetNotFound(t *testing.T) {
	s := setupTestStore(t)

	_, err := s.Get("US-999")
	if err == nil {
		t.Fatal("expected error for missing task")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error = %q, want it to contain 'not found'", err.Error())
	}
}

func TestUpdate(t *testing.T) {
	s := setupTestStore(t)

	tk := &task.Task{
		ID:      "US-001",
		Title:   "Test task",
		Status:  "backlog",
		Created: "2026-01-01",
		Updated: "2026-01-01",
	}
	s.Create(tk)

	tk.Status = "in-progress"
	tk.Assignee = "claude"
	if err := s.Update(tk); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, _ := s.Get("US-001")
	if got.Status != "in-progress" {
		t.Errorf("Status = %q, want %q", got.Status, "in-progress")
	}
	if got.Assignee != "claude" {
		t.Errorf("Assignee = %q, want %q", got.Assignee, "claude")
	}
}

func TestList(t *testing.T) {
	s := setupTestStore(t)

	tasks := []*task.Task{
		{ID: "US-001", Title: "Task 1", Status: "backlog", Priority: "high", Created: "2026-01-01", Updated: "2026-01-01"},
		{ID: "US-002", Title: "Task 2", Status: "in-progress", Priority: "low", Assignee: "andy", Created: "2026-01-01", Updated: "2026-01-01"},
		{ID: "US-003", Title: "Task 3", Status: "backlog", Priority: "high", Tags: task.FlowSlice{"api"}, Created: "2026-01-01", Updated: "2026-01-01"},
	}
	for _, tk := range tasks {
		s.Create(tk)
	}

	// No filter
	all, _ := s.List(Filter{})
	if len(all) != 3 {
		t.Errorf("List() = %d tasks, want 3", len(all))
	}

	// Filter by status
	backlog, _ := s.List(Filter{Status: "backlog"})
	if len(backlog) != 2 {
		t.Errorf("List(status=backlog) = %d tasks, want 2", len(backlog))
	}

	// Filter by assignee
	assigned, _ := s.List(Filter{Assignee: "andy"})
	if len(assigned) != 1 {
		t.Errorf("List(assignee=andy) = %d tasks, want 1", len(assigned))
	}

	// Filter by tag
	tagged, _ := s.List(Filter{Tag: "api"})
	if len(tagged) != 1 {
		t.Errorf("List(tag=api) = %d tasks, want 1", len(tagged))
	}

	// Combined filter
	combo, _ := s.List(Filter{Status: "backlog", Priority: "high"})
	if len(combo) != 2 {
		t.Errorf("List(status=backlog,priority=high) = %d tasks, want 2", len(combo))
	}
}

func TestNextID(t *testing.T) {
	s := setupTestStore(t)

	// First ID
	id1, err := s.NextID()
	if err != nil {
		t.Fatalf("NextID: %v", err)
	}
	if id1 != "US-001" {
		t.Errorf("first NextID = %q, want %q", id1, "US-001")
	}

	// Create a task and get the next
	s.Create(&task.Task{ID: id1, Title: "First", Status: "backlog", Created: "2026-01-01", Updated: "2026-01-01"})
	id2, _ := s.NextID()
	if id2 != "US-002" {
		t.Errorf("second NextID = %q, want %q", id2, "US-002")
	}
}

func TestLoadTemplate(t *testing.T) {
	s := setupTestStore(t)

	body, err := s.LoadTemplate("default")
	if err != nil {
		t.Fatalf("LoadTemplate(default): %v", err)
	}
	if !strings.Contains(body, "Acceptance Criteria") {
		t.Error("default template missing 'Acceptance Criteria'")
	}

	_, err = s.LoadTemplate("nonexistent")
	if err == nil {
		t.Error("expected error for missing template")
	}
}

func TestAutoCommitCreate(t *testing.T) {
	s, dir := setupTestStoreWithGit(t)
	s.Config.AutoCommit = true

	tk := &task.Task{
		ID:      "US-001",
		Title:   "Auto-commit test",
		Status:  "backlog",
		Created: "2026-01-01",
		Updated: "2026-01-01",
	}

	if err := s.Create(tk); err != nil {
		t.Fatalf("Create: %v", err)
	}

	commits := gitLog(t, dir)
	if len(commits) < 2 {
		t.Fatalf("expected at least 2 commits, got %d", len(commits))
	}

	latest := commits[0]
	if !strings.Contains(latest, "skeeter: create US-001") {
		t.Errorf("commit message = %q, want it to contain 'skeeter: create US-001'", latest)
	}
}

func TestAutoCommitUpdate(t *testing.T) {
	s, dir := setupTestStoreWithGit(t)
	s.Config.AutoCommit = true

	tk := &task.Task{
		ID:      "US-001",
		Title:   "Auto-commit test",
		Status:  "backlog",
		Created: "2026-01-01",
		Updated: "2026-01-01",
	}
	s.Create(tk)

	tk.Status = "in-progress"
	if err := s.Update(tk); err != nil {
		t.Fatalf("Update: %v", err)
	}

	commits := gitLog(t, dir)
	latest := commits[0]
	if !strings.Contains(latest, "skeeter: update US-001") {
		t.Errorf("commit message = %q, want it to contain 'skeeter: update US-001'", latest)
	}
}

func TestAutoCommitDisabled(t *testing.T) {
	s, dir := setupTestStoreWithGit(t)
	// auto_commit defaults to false

	tk := &task.Task{
		ID:      "US-001",
		Title:   "No auto-commit",
		Status:  "backlog",
		Created: "2026-01-01",
		Updated: "2026-01-01",
	}
	s.Create(tk)

	commits := gitLog(t, dir)
	// Should only have the initial commit
	if len(commits) != 1 {
		t.Errorf("expected 1 commit (initial only), got %d: %v", len(commits), commits)
	}
}
