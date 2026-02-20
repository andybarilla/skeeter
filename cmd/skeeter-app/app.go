package main

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/andybarilla/skeeter/internal/config"
	"github.com/andybarilla/skeeter/internal/store"
	"github.com/andybarilla/skeeter/internal/task"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type BoardData struct {
	Columns  []ColumnData   `json:"columns"`
	Config   *config.Config `json:"config"`
	RepoName string         `json:"repoName"`
}

type ColumnData struct {
	Status string      `json:"status"`
	Tasks  []task.Task `json:"tasks"`
}

type BoardFilter struct {
	Priority string `json:"priority"`
	Assignee string `json:"assignee"`
	Tag      string `json:"tag"`
}

type CreateTaskInput struct {
	Title    string   `json:"title"`
	Priority string   `json:"priority"`
	Assignee string   `json:"assignee"`
	Tags     []string `json:"tags"`
	Body     string   `json:"body"`
}

type UpdateTaskInput struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Status   string   `json:"status"`
	Priority string   `json:"priority"`
	Assignee string   `json:"assignee"`
	Tags     []string `json:"tags"`
	Body     string   `json:"body"`
}

type App struct {
	ctx       context.Context
	mu        sync.RWMutex
	store     store.Store
	repoStore *RepoStore
	repoName  string
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	rs, err := NewRepoStore()
	if err != nil {
		return
	}
	a.repoStore = rs

	// Auto-open the first (MRU) repo
	repos, err := rs.Load()
	if err != nil || len(repos) == 0 {
		return
	}
	s, err := openStoreFromEntry(repos[0])
	if err != nil {
		return
	}
	a.mu.Lock()
	a.store = s
	a.repoName = repos[0].Name
	a.mu.Unlock()
}

func openStoreFromEntry(entry RepoEntry) (store.Store, error) {
	if entry.Remote != "" {
		return store.NewGitHub(entry.Remote, entry.Dir)
	}
	return store.NewFilesystem(entry.Path)
}

// GetBoard returns all tasks grouped by status columns.
func (a *App) GetBoard(filter BoardFilter) (*BoardData, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.store == nil {
		return &BoardData{}, nil
	}

	cfg := a.store.GetConfig()

	// List all tasks (no status filter — we need all columns)
	f := store.Filter{
		Priority: filter.Priority,
		Assignee: filter.Assignee,
		Tag:      filter.Tag,
	}
	tasks, err := a.store.List(f)
	if err != nil {
		return nil, err
	}

	// Group by status
	grouped := make(map[string][]task.Task)
	for _, t := range tasks {
		grouped[t.Status] = append(grouped[t.Status], t)
	}

	// Sort each group by priority rank then ID
	for status := range grouped {
		g := grouped[status]
		sort.Slice(g, func(i, j int) bool {
			ri := cfg.PriorityRank(g[i].Priority)
			rj := cfg.PriorityRank(g[j].Priority)
			if ri != rj {
				return ri < rj
			}
			return g[i].ID < g[j].ID
		})
		grouped[status] = g
	}

	// Build columns in configured order
	var columns []ColumnData
	for _, status := range cfg.Statuses {
		columns = append(columns, ColumnData{
			Status: status,
			Tasks:  grouped[status],
		})
	}

	return &BoardData{
		Columns:  columns,
		Config:   cfg,
		RepoName: a.repoName,
	}, nil
}

// GetTask returns a single task by ID.
func (a *App) GetTask(id string) (*task.Task, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.store == nil {
		return nil, fmt.Errorf("no repo selected")
	}
	return a.store.Get(id)
}

// CreateTask creates a new task.
func (a *App) CreateTask(input CreateTaskInput) (*task.Task, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.store == nil {
		return nil, fmt.Errorf("no repo selected")
	}

	cfg := a.store.GetConfig()
	id, err := a.store.NextID()
	if err != nil {
		return nil, err
	}

	now := time.Now().Format("2006-01-02")
	priority := input.Priority
	if priority == "" {
		priority = cfg.Priorities[len(cfg.Priorities)-1] // lowest
	}

	t := &task.Task{
		ID:       id,
		Title:    input.Title,
		Status:   cfg.Statuses[0], // first status (backlog)
		Priority: priority,
		Assignee: input.Assignee,
		Tags:     input.Tags,
		Created:  now,
		Updated:  now,
		Body:     input.Body,
	}

	if err := a.store.Create(t); err != nil {
		return nil, err
	}
	return t, nil
}

// UpdateTask performs a full task update.
func (a *App) UpdateTask(input UpdateTaskInput) (*task.Task, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.store == nil {
		return nil, fmt.Errorf("no repo selected")
	}

	t, err := a.store.Get(input.ID)
	if err != nil {
		return nil, err
	}

	t.Title = input.Title
	t.Status = input.Status
	t.Priority = input.Priority
	t.Assignee = input.Assignee
	t.Tags = input.Tags
	t.Body = input.Body
	t.Updated = time.Now().Format("2006-01-02")

	if err := a.store.Update(t); err != nil {
		return nil, err
	}
	return t, nil
}

// MoveTask changes only the status of a task (for drag-and-drop).
func (a *App) MoveTask(id, status string) (*task.Task, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.store == nil {
		return nil, fmt.Errorf("no repo selected")
	}

	cfg := a.store.GetConfig()
	if !cfg.ValidStatus(status) {
		return nil, fmt.Errorf("invalid status %q", status)
	}

	t, err := a.store.Get(id)
	if err != nil {
		return nil, err
	}

	t.Status = status
	t.Updated = time.Now().Format("2006-01-02")

	if err := a.store.Update(t); err != nil {
		return nil, err
	}
	return t, nil
}

// AssignTask changes only the assignee of a task.
func (a *App) AssignTask(id, assignee string) (*task.Task, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.store == nil {
		return nil, fmt.Errorf("no repo selected")
	}

	t, err := a.store.Get(id)
	if err != nil {
		return nil, err
	}

	t.Assignee = assignee
	t.Updated = time.Now().Format("2006-01-02")

	if err := a.store.Update(t); err != nil {
		return nil, err
	}
	return t, nil
}

// GetRepos returns the saved repo list.
func (a *App) GetRepos() ([]RepoEntry, error) {
	return a.repoStore.Load()
}

// AddRepo validates and saves a new repo.
func (a *App) AddRepo(entry RepoEntry) error {
	// Validate by opening the store
	s, err := openStoreFromEntry(entry)
	if err != nil {
		return fmt.Errorf("cannot open repo: %w", err)
	}

	// Use project name as display name if not provided
	if entry.Name == "" {
		cfg := s.GetConfig()
		entry.Name = cfg.Project.Name
		if entry.Name == "" {
			entry.Name = "unnamed"
		}
	}

	if err := a.repoStore.Add(entry); err != nil {
		return err
	}

	// If no active store, switch to the new one
	a.mu.RLock()
	hasStore := a.store != nil
	a.mu.RUnlock()

	if !hasStore {
		a.mu.Lock()
		a.store = s
		a.repoName = entry.Name
		a.mu.Unlock()
		runtime.WindowSetTitle(a.ctx, "Skeeter — "+entry.Name)
	}

	return nil
}

// RemoveRepo removes a saved repo.
func (a *App) RemoveRepo(name string) error {
	// If removing the active repo, clear the store
	a.mu.Lock()
	if a.repoName == name {
		a.store = nil
		a.repoName = ""
	}
	a.mu.Unlock()

	return a.repoStore.Remove(name)
}

// SwitchRepo switches the active store.
func (a *App) SwitchRepo(name string) error {
	repos, err := a.repoStore.Load()
	if err != nil {
		return err
	}

	var entry *RepoEntry
	for _, r := range repos {
		if r.Name == name {
			r := r
			entry = &r
			break
		}
	}
	if entry == nil {
		return fmt.Errorf("repo %q not found", name)
	}

	s, err := openStoreFromEntry(*entry)
	if err != nil {
		return fmt.Errorf("cannot open repo: %w", err)
	}

	a.mu.Lock()
	a.store = s
	a.repoName = name
	a.mu.Unlock()

	// Move to front (MRU)
	_ = a.repoStore.MoveToFront(name)

	runtime.WindowSetTitle(a.ctx, "Skeeter — "+name)
	return nil
}

// GetActiveRepoName returns the current repo name.
func (a *App) GetActiveRepoName() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.repoName
}
