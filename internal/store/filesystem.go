package store

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/andybarilla/skeeter/internal/config"
	"github.com/andybarilla/skeeter/internal/id"
	"github.com/andybarilla/skeeter/internal/task"
)

type FilesystemStore struct {
	Dir    string
	Config *config.Config
}

func NewFilesystem(dir string) (*FilesystemStore, error) {
	cfg, err := config.Load(dir)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	return &FilesystemStore{Dir: dir, Config: cfg}, nil
}

func (s *FilesystemStore) tasksDir() string {
	return filepath.Join(s.Dir, "tasks")
}

func (s *FilesystemStore) taskPath(taskID string) string {
	return filepath.Join(s.tasksDir(), taskID+".md")
}

func (s *FilesystemStore) Init(projectName string) error {
	if err := os.MkdirAll(s.tasksDir(), 0755); err != nil {
		return err
	}

	cfg := config.Default()
	cfg.Project.Name = projectName
	if err := cfg.Save(s.Dir); err != nil {
		return err
	}

	s.Config = cfg
	return s.writeSkeeterMD()
}

func (s *FilesystemStore) writeSkeeterMD() error {
	content := `# Skeeter — Project Tasks

Tasks are markdown files with YAML frontmatter in the ` + "`tasks/`" + ` subdirectory.

## For Agents: Finding Work

1. Look for tasks where ` + "`status: ready-for-development`" + ` and ` + "`assignee:`" + ` is empty
2. Set ` + "`assignee: <your-name>`" + ` and ` + "`status: in-progress`" + ` before starting
3. Use ` + "`Acceptance Criteria`" + ` as your definition of done
4. Set ` + "`status: done`" + ` when complete

## Frontmatter Fields

| Field      | Description                                              |
|------------|----------------------------------------------------------|
| id         | Task identifier (e.g., US-001)                           |
| title      | Short task title                                         |
| status     | One of: ` + strings.Join(s.Config.Statuses, ", ") + `   |
| priority   | One of: ` + strings.Join(s.Config.Priorities, ", ") + `  |
| assignee   | Who is working on this (empty = available)               |
| tags       | Array of labels                                          |
| links      | Related URLs                                             |
| created    | Creation date                                            |
| updated    | Last modified date                                       |
`
	path := filepath.Join(s.Dir, "SKEETER.md")
	return os.WriteFile(path, []byte(content), 0644)
}

func (s *FilesystemStore) List(filter Filter) ([]task.Task, error) {
	entries, err := os.ReadDir(s.tasksDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var tasks []task.Task
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".md" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(s.tasksDir(), e.Name()))
		if err != nil {
			continue
		}

		t, err := task.Parse(string(data))
		if err != nil {
			continue
		}

		if !matchesFilter(t, filter) {
			continue
		}

		tasks = append(tasks, *t)
	}
	return tasks, nil
}

func matchesFilter(t *task.Task, f Filter) bool {
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.Priority != "" && t.Priority != f.Priority {
		return false
	}
	if f.Assignee != "" && t.Assignee != f.Assignee {
		return false
	}
	if f.Tag != "" {
		found := false
		for _, tag := range t.Tags {
			if tag == f.Tag {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (s *FilesystemStore) Get(taskID string) (*task.Task, error) {
	data, err := os.ReadFile(s.taskPath(taskID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("task %s not found", taskID)
		}
		return nil, err
	}
	return task.Parse(string(data))
}

func (s *FilesystemStore) Create(t *task.Task) error {
	content, err := task.Marshal(t)
	if err != nil {
		return err
	}
	path := s.taskPath(t.ID)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}
	return s.autoCommit(fmt.Sprintf("create %s: %s", t.ID, t.Title), path)
}

func (s *FilesystemStore) Update(t *task.Task) error {
	t.Updated = time.Now().Format("2006-01-02")
	content, err := task.Marshal(t)
	if err != nil {
		return err
	}
	path := s.taskPath(t.ID)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}
	return s.autoCommit(fmt.Sprintf("update %s: %s", t.ID, t.Title), path)
}

func (s *FilesystemStore) NextID() (string, error) {
	return id.Next(s.tasksDir(), s.Config.Project.Prefix)
}

func (s *FilesystemStore) GetConfig() *config.Config {
	return s.Config
}

func (s *FilesystemStore) autoCommit(message string, files ...string) error {
	if !s.Config.AutoCommit {
		return nil
	}

	args := append([]string{"add"}, files...)
	if err := exec.Command("git", args...).Run(); err != nil {
		return fmt.Errorf("git add: %w", err)
	}

	if err := exec.Command("git", "commit", "-m", "skeeter: "+message).Run(); err != nil {
		return fmt.Errorf("git commit: %w", err)
	}
	return nil
}
