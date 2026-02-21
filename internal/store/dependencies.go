package store

import (
	"github.com/andybarilla/skeeter/internal/task"
)

type DependencyStatus struct {
	DependsOn          []string
	BlockedBy          []string
	Blocking           []string
	AllDependenciesMet bool
}

func GetDependencyStatus(s Store, t *task.Task, allTasks []task.Task) *DependencyStatus {
	status := &DependencyStatus{
		AllDependenciesMet: true,
	}

	doneStatus := getDoneStatus(s)

	taskMap := make(map[string]*task.Task)
	for i := range allTasks {
		taskMap[allTasks[i].ID] = &allTasks[i]
	}

	for _, depID := range t.DependsOn {
		status.DependsOn = append(status.DependsOn, depID)
		depTask, exists := taskMap[depID]
		if !exists || depTask.Status != doneStatus {
			status.BlockedBy = append(status.BlockedBy, depID)
			status.AllDependenciesMet = false
		}
	}

	for i := range allTasks {
		for _, depID := range allTasks[i].DependsOn {
			if depID == t.ID {
				status.Blocking = append(status.Blocking, allTasks[i].ID)
			}
		}
	}

	return status
}

func IsBlocked(s Store, t *task.Task, allTasks []task.Task) bool {
	if len(t.DependsOn) == 0 {
		return false
	}

	doneStatus := getDoneStatus(s)

	taskMap := make(map[string]string)
	for _, tt := range allTasks {
		taskMap[tt.ID] = tt.Status
	}

	for _, depID := range t.DependsOn {
		status, exists := taskMap[depID]
		if !exists || status != doneStatus {
			return true
		}
	}

	return false
}

func getDoneStatus(s Store) string {
	cfg := s.GetConfig()
	if len(cfg.Statuses) > 0 {
		return cfg.Statuses[len(cfg.Statuses)-1]
	}
	return "done"
}

func DetectCircularDependency(t *task.Task, s Store) ([]string, error) {
	visited := make(map[string]bool)
	path := []string{}
	return detectCycle(t.ID, t.DependsOn, s, visited, path)
}

func detectCycle(startID string, dependsOn []string, s Store, visited map[string]bool, path []string) ([]string, error) {
	for _, depID := range dependsOn {
		if depID == startID {
			return append(path, depID), nil
		}

		if visited[depID] {
			continue
		}
		visited[depID] = true

		depTask, err := s.Get(depID)
		if err != nil {
			continue
		}

		newPath := append(path, depID)
		if cycle, err := detectCycle(startID, depTask.DependsOn, s, visited, newPath); err == nil && len(cycle) > 0 {
			return cycle, nil
		}
	}

	return nil, nil
}
