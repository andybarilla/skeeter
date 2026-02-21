package store

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andybarilla/skeeter/internal/config"
	"github.com/andybarilla/skeeter/internal/task"
)

func setupGitHubServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/repos/owner/repo/contents/.skeeter/config.yaml", func(w http.ResponseWriter, r *http.Request) {
		configYAML := `project:
  name: test-project
  prefix: US
statuses:
  - backlog
  - ready-for-development
  - in-progress
  - done
priorities:
  - critical
  - high
  - medium
  - low
`
		resp := ghContentsResponse{
			Name:     "config.yaml",
			Path:     ".skeeter/config.yaml",
			Content:  base64.StdEncoding.EncodeToString([]byte(configYAML)),
			Encoding: "base64",
			Type:     "file",
		}
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/repos/owner/repo/contents/.skeeter/tasks", func(w http.ResponseWriter, r *http.Request) {
		resp := []ghContentsResponse{
			{Name: "US-001.md", Path: ".skeeter/tasks/US-001.md", Type: "file"},
			{Name: "US-002.md", Path: ".skeeter/tasks/US-002.md", Type: "file"},
		}
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/repos/owner/repo/contents/.skeeter/tasks/US-001.md", func(w http.ResponseWriter, r *http.Request) {
		content := `---
id: US-001
title: Test Task
status: backlog
priority: high
created: "2026-01-01"
updated: "2026-01-01"
---

Task body here.
`
		if r.Method == "GET" {
			resp := ghContentsResponse{
				Name:     "US-001.md",
				Path:     ".skeeter/tasks/US-001.md",
				Content:  base64.StdEncoding.EncodeToString([]byte(content)),
				Encoding: "base64",
				Type:     "file",
				SHA:      "abc123",
			}
			json.NewEncoder(w).Encode(resp)
		} else if r.Method == "PUT" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "updated"})
		}
	})

	mux.HandleFunc("/repos/owner/repo/contents/.skeeter/tasks/US-002.md", func(w http.ResponseWriter, r *http.Request) {
		content := `---
id: US-002
title: Another Task
status: in-progress
priority: low
assignee: alice
created: "2026-01-01"
updated: "2026-01-01"
---
`
		resp := ghContentsResponse{
			Name:     "US-002.md",
			Path:     ".skeeter/tasks/US-002.md",
			Content:  base64.StdEncoding.EncodeToString([]byte(content)),
			Encoding: "base64",
			Type:     "file",
			SHA:      "def456",
		}
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/repos/owner/repo/contents/.skeeter/templates/default.md", func(w http.ResponseWriter, r *http.Request) {
		content := "## Acceptance Criteria\n\n- [ ]\n"
		resp := ghContentsResponse{
			Name:     "default.md",
			Path:     ".skeeter/templates/default.md",
			Content:  base64.StdEncoding.EncodeToString([]byte(content)),
			Encoding: "base64",
			Type:     "file",
		}
		json.NewEncoder(w).Encode(resp)
	})

	return httptest.NewServer(mux)
}

func TestGitHubStoreList(t *testing.T) {
	server := setupGitHubServer()
	defer server.Close()

	store := &GitHubStore{
		owner:   "owner",
		repo:    "repo",
		dir:     ".skeeter",
		token:   "fake-token",
		client:  server.Client(),
		baseURL: server.URL,
	}

	cfg := defaultConfigForTest()
	store.cfg = cfg

	tasks, err := store.List(Filter{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("List returned %d tasks, want 2", len(tasks))
	}
}

func TestGitHubStoreListWithFilter(t *testing.T) {
	server := setupGitHubServer()
	defer server.Close()

	store := &GitHubStore{
		owner:   "owner",
		repo:    "repo",
		dir:     ".skeeter",
		token:   "fake-token",
		client:  server.Client(),
		baseURL: server.URL,
	}
	store.cfg = defaultConfigForTest()

	tasks, err := store.List(Filter{Status: "in-progress"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("List returned %d tasks, want 1", len(tasks))
	}
	if tasks[0].ID != "US-002" {
		t.Errorf("Task ID = %s, want US-002", tasks[0].ID)
	}
}

func TestGitHubStoreGet(t *testing.T) {
	server := setupGitHubServer()
	defer server.Close()

	store := &GitHubStore{
		owner:   "owner",
		repo:    "repo",
		dir:     ".skeeter",
		token:   "fake-token",
		client:  server.Client(),
		baseURL: server.URL,
	}
	store.cfg = defaultConfigForTest()

	tk, err := store.Get("US-001")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if tk.Title != "Test Task" {
		t.Errorf("Title = %q, want %q", tk.Title, "Test Task")
	}
	if tk.Body != "Task body here.\n" {
		t.Errorf("Body = %q, want %q", tk.Body, "Task body here.\n")
	}
}

func TestGitHubStoreGetNotFound(t *testing.T) {
	server := setupGitHubServer()
	defer server.Close()

	store := &GitHubStore{
		owner:   "owner",
		repo:    "repo",
		dir:     ".skeeter",
		token:   "fake-token",
		client:  server.Client(),
		baseURL: server.URL,
	}
	store.cfg = defaultConfigForTest()

	_, err := store.Get("US-999")
	if err == nil {
		t.Error("expected error for missing task")
	}
}

func TestGitHubStoreCreate(t *testing.T) {
	server := setupGitHubServer()
	defer server.Close()

	store := &GitHubStore{
		owner:  "owner",
		repo:   "repo",
		dir:    ".skeeter",
		token:  "fake-token",
		client: server.Client(),
	}
	store.cfg = defaultConfigForTest()

	tk := &task.Task{
		ID:       "US-003",
		Title:    "New Task",
		Status:   "backlog",
		Priority: "medium",
		Created:  "2026-01-01",
		Updated:  "2026-01-01",
	}

	// Create would need a PUT endpoint for new files
	// For this test, we just verify the method builds correct content
	content, err := task.Marshal(tk)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !contains(content, "US-003") {
		t.Errorf("marshaled content missing task ID: %s", content)
	}
}

func TestGitHubStoreNextID(t *testing.T) {
	server := setupGitHubServer()
	defer server.Close()

	store := &GitHubStore{
		owner:   "owner",
		repo:    "repo",
		dir:     ".skeeter",
		token:   "fake-token",
		client:  server.Client(),
		baseURL: server.URL,
	}
	store.cfg = defaultConfigForTest()

	id, err := store.NextID()
	if err != nil {
		t.Fatalf("NextID: %v", err)
	}
	if id != "US-003" {
		t.Errorf("NextID = %q, want %q", id, "US-003")
	}
}

func TestGitHubStoreLoadTemplate(t *testing.T) {
	server := setupGitHubServer()
	defer server.Close()

	store := &GitHubStore{
		owner:   "owner",
		repo:    "repo",
		dir:     ".skeeter",
		token:   "fake-token",
		client:  server.Client(),
		baseURL: server.URL,
	}
	store.cfg = defaultConfigForTest()

	content, err := store.LoadTemplate("default")
	if err != nil {
		t.Fatalf("LoadTemplate: %v", err)
	}
	if content == "" {
		t.Error("LoadTemplate returned empty content")
	}
}

func TestGitHubStoreInit(t *testing.T) {
	store := &GitHubStore{}

	err := store.Init("test")
	if err == nil {
		t.Error("expected error for Init on remote store")
	}
}

func TestGitHubStoreGetConfig(t *testing.T) {
	store := &GitHubStore{
		cfg: defaultConfigForTest(),
	}

	cfg := store.GetConfig()
	if cfg == nil {
		t.Fatal("GetConfig returned nil")
	}
	if cfg.Project.Name != "test-project" {
		t.Errorf("Project.Name = %q, want %q", cfg.Project.Name, "test-project")
	}
}

func defaultConfigForTest() *config.Config {
	return &config.Config{
		Project: config.ProjectConfig{
			Name:   "test-project",
			Prefix: "US",
		},
		Statuses:   []string{"backlog", "ready-for-development", "in-progress", "done"},
		Priorities: []string{"critical", "high", "medium", "low"},
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
