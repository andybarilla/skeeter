package store

import (
	"fmt"
	"testing"

	"github.com/andybarilla/skeeter/internal/config"
	"github.com/andybarilla/skeeter/internal/task"
)

type mockStore struct {
	tasks   map[string]*task.Task
	config  *config.Config
	taskIDs []string
}

func newMockStore() *mockStore {
	return &mockStore{
		tasks:  make(map[string]*task.Task),
		config: config.Default(),
	}
}

func (m *mockStore) Init(projectName string) error { return nil }
func (m *mockStore) List(filter Filter) ([]task.Task, error) {
	var result []task.Task
	for _, id := range m.taskIDs {
		if t, ok := m.tasks[id]; ok {
			if filter.Status != "" && t.Status != filter.Status {
				continue
			}
			result = append(result, *t)
		}
	}
	return result, nil
}
func (m *mockStore) Get(id string) (*task.Task, error) {
	if t, ok := m.tasks[id]; ok {
		return t, nil
	}
	return nil, ErrTaskNotFound
}
func (m *mockStore) Create(t *task.Task) error {
	m.tasks[t.ID] = t
	m.taskIDs = append(m.taskIDs, t.ID)
	return nil
}
func (m *mockStore) Update(t *task.Task) error {
	m.tasks[t.ID] = t
	return nil
}
func (m *mockStore) NextID() (string, error) {
	return "US-999", nil
}
func (m *mockStore) GetConfig() *config.Config { return m.config }
func (m *mockStore) LoadTemplate(name string) (string, error) {
	return "template", nil
}

var ErrTaskNotFound = fmt.Errorf("task not found")

func TestIsBlocked(t *testing.T) {
	s := newMockStore()

	s.Create(&task.Task{ID: "US-001", Status: "done"})
	s.Create(&task.Task{ID: "US-002", Status: "in-progress"})
	s.Create(&task.Task{ID: "US-003", Status: "backlog", DependsOn: task.FlowSlice{"US-001", "US-002"}})
	s.Create(&task.Task{ID: "US-004", Status: "backlog", DependsOn: task.FlowSlice{"US-001"}})
	s.Create(&task.Task{ID: "US-005", Status: "backlog"})

	allTasks, _ := s.List(Filter{})

	t.Run("not blocked with no dependencies", func(t *testing.T) {
		if IsBlocked(s, &task.Task{ID: "US-005"}, allTasks) {
			t.Error("task with no dependencies should not be blocked")
		}
	})

	t.Run("not blocked when all dependencies complete", func(t *testing.T) {
		tk := &task.Task{ID: "US-004", DependsOn: task.FlowSlice{"US-001"}}
		if IsBlocked(s, tk, allTasks) {
			t.Error("task with complete dependencies should not be blocked")
		}
	})

	t.Run("blocked when dependency incomplete", func(t *testing.T) {
		tk := &task.Task{ID: "US-003", DependsOn: task.FlowSlice{"US-001", "US-002"}}
		if !IsBlocked(s, tk, allTasks) {
			t.Error("task with incomplete dependencies should be blocked")
		}
	})

	t.Run("blocked when dependency missing", func(t *testing.T) {
		tk := &task.Task{ID: "US-006", DependsOn: task.FlowSlice{"US-999"}}
		if !IsBlocked(s, tk, allTasks) {
			t.Error("task with missing dependency should be blocked")
		}
	})
}

func TestGetDependencyStatus(t *testing.T) {
	s := newMockStore()

	s.Create(&task.Task{ID: "US-001", Status: "done"})
	s.Create(&task.Task{ID: "US-002", Status: "in-progress"})
	s.Create(&task.Task{ID: "US-003", Status: "backlog", DependsOn: task.FlowSlice{"US-001", "US-002"}})
	s.Create(&task.Task{ID: "US-004", Status: "backlog", DependsOn: task.FlowSlice{"US-003"}})

	allTasks, _ := s.List(Filter{})

	t.Run("shows blocked by", func(t *testing.T) {
		tk := &task.Task{ID: "US-003", DependsOn: task.FlowSlice{"US-001", "US-002"}}
		status := GetDependencyStatus(s, tk, allTasks)

		if len(status.DependsOn) != 2 {
			t.Errorf("DependsOn = %v, want 2 items", status.DependsOn)
		}
		if len(status.BlockedBy) != 1 {
			t.Errorf("BlockedBy = %v, want 1 item", status.BlockedBy)
		}
		if status.AllDependenciesMet {
			t.Error("AllDependenciesMet should be false")
		}
	})

	t.Run("shows blocking", func(t *testing.T) {
		tk := &task.Task{ID: "US-001"}
		status := GetDependencyStatus(s, tk, allTasks)

		if len(status.Blocking) != 1 {
			t.Errorf("Blocking = %v, want 1 item", status.Blocking)
		}
		if status.Blocking[0] != "US-003" {
			t.Errorf("Blocking = %v, want [US-003]", status.Blocking)
		}
	})

	t.Run("all dependencies met", func(t *testing.T) {
		tk := &task.Task{ID: "US-004", DependsOn: task.FlowSlice{"US-001"}}
		status := GetDependencyStatus(s, tk, allTasks)

		if !status.AllDependenciesMet {
			t.Error("AllDependenciesMet should be true when only dependency is done")
		}
	})
}

func TestDetectCircularDependency(t *testing.T) {
	t.Run("detects direct self-reference", func(t *testing.T) {
		s := newMockStore()
		s.Create(&task.Task{ID: "US-001", DependsOn: task.FlowSlice{"US-002"}})

		tk := &task.Task{ID: "US-002", DependsOn: task.FlowSlice{"US-002"}}
		cycle, _ := DetectCircularDependency(tk, s)
		if len(cycle) == 0 {
			t.Error("should detect self-referencing dependency")
		}
	})

	t.Run("detects cycle back to new task", func(t *testing.T) {
		s := newMockStore()
		s.Create(&task.Task{ID: "US-001", DependsOn: task.FlowSlice{"US-002"}})

		tk := &task.Task{ID: "US-002", DependsOn: task.FlowSlice{"US-001"}}
		cycle, _ := DetectCircularDependency(tk, s)
		if len(cycle) == 0 {
			t.Error("should detect cycle: US-002 -> US-001 -> US-002")
		}
	})

	t.Run("no cycle with linear chain", func(t *testing.T) {
		s := newMockStore()
		s.Create(&task.Task{ID: "US-001"})
		s.Create(&task.Task{ID: "US-002", DependsOn: task.FlowSlice{"US-001"}})

		tk := &task.Task{ID: "US-003", DependsOn: task.FlowSlice{"US-002"}}
		cycle, _ := DetectCircularDependency(tk, s)
		if len(cycle) > 0 {
			t.Error("should not detect cycle in linear chain")
		}
	})

	t.Run("detects longer cycle", func(t *testing.T) {
		s := newMockStore()
		s.Create(&task.Task{ID: "US-001", DependsOn: task.FlowSlice{"US-003"}})
		s.Create(&task.Task{ID: "US-002", DependsOn: task.FlowSlice{"US-001"}})

		tk := &task.Task{ID: "US-003", DependsOn: task.FlowSlice{"US-002"}}
		cycle, _ := DetectCircularDependency(tk, s)
		if len(cycle) == 0 {
			t.Error("should detect cycle: US-003 -> US-002 -> US-001 -> US-003")
		}
	})
}
