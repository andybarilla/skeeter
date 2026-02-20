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

func (s *FilesystemStore) templatesDir() string {
	return filepath.Join(s.Dir, "templates")
}

func (s *FilesystemStore) Init(projectName string) error {
	if err := os.MkdirAll(s.tasksDir(), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(s.templatesDir(), 0755); err != nil {
		return err
	}

	cfg := config.Default()
	cfg.Project.Name = projectName
	if err := cfg.Save(s.Dir); err != nil {
		return err
	}

	s.Config = cfg

	if err := s.writeDefaultTemplate(); err != nil {
		return err
	}

	return s.writeSkeeterMD()
}

func (s *FilesystemStore) writeDefaultTemplate() error {
	content := `## Acceptance Criteria

- [ ]

## Context

`
	path := filepath.Join(s.templatesDir(), "default.md")
	return os.WriteFile(path, []byte(content), 0644)
}

func (s *FilesystemStore) writeSkeeterMD() error {
	// Derive the "ready" status (second in list) and "done" status (last in list)
	readyStatus := "ready-for-development"
	inProgressStatus := "in-progress"
	doneStatus := "done"
	if len(s.Config.Statuses) >= 2 {
		readyStatus = s.Config.Statuses[1]
	}
	if len(s.Config.Statuses) >= 3 {
		inProgressStatus = s.Config.Statuses[2]
	}
	if len(s.Config.Statuses) >= 1 {
		doneStatus = s.Config.Statuses[len(s.Config.Statuses)-1]
	}

	prefix := s.Config.Project.Prefix

	content := "# Skeeter — Project Tasks\n\n" +
		"Tasks are markdown files with YAML frontmatter in the `tasks/` subdirectory.\n\n" +
		"## For Agents: Finding Work\n\n" +
		"1. Look for tasks where `status: " + readyStatus + "` and `assignee:` is empty\n" +
		"2. Set `assignee: <your-name>` and `status: " + inProgressStatus + "` before starting\n" +
		"3. Use `Acceptance Criteria` as your definition of done\n" +
		"4. Set `status: " + doneStatus + "` when complete\n\n" +
		"## Frontmatter Fields\n\n" +
		"| Field      | Description                                              |\n" +
		"|------------|----------------------------------------------------------|\n" +
		"| id         | Task identifier (e.g., " + prefix + "-001)                           |\n" +
		"| title      | Short task title                                         |\n" +
		"| status     | One of: " + strings.Join(s.Config.Statuses, ", ") + " |\n" +
		"| priority   | One of: " + strings.Join(s.Config.Priorities, ", ") + " |\n" +
		"| assignee   | Who is working on this (empty = available)               |\n" +
		"| tags       | Array of labels                                          |\n" +
		"| links      | Related URLs                                             |\n" +
		"| created    | Creation date                                            |\n" +
		"| updated    | Last modified date                                       |\n"

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
	s.writeSkeeterMD()
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
	s.writeSkeeterMD()
	return s.autoCommit(fmt.Sprintf("update %s: %s", t.ID, t.Title), path)
}

func (s *FilesystemStore) NextID() (string, error) {
	return id.Next(s.tasksDir(), s.Config.Project.Prefix)
}

func (s *FilesystemStore) GetConfig() *config.Config {
	return s.Config
}

func (s *FilesystemStore) LoadTemplate(name string) (string, error) {
	path := filepath.Join(s.templatesDir(), name+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("template %q not found (looked in %s)", name, s.templatesDir())
		}
		return "", err
	}
	return string(data), nil
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
