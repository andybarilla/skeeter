package main

import (
	"sort"

	"github.com/andybarilla/skeeter/internal/config"
	"github.com/andybarilla/skeeter/internal/store"
	"github.com/andybarilla/skeeter/internal/task"
)

// pickNextTask returns the highest-priority unassigned task in ready-for-development status.
// Returns nil with no error when no tasks are available.
func pickNextTask(s store.Store, cfg *config.Config) (*task.Task, error) {
	tasks, err := s.List(store.Filter{Status: "ready-for-development"})
	if err != nil {
		return nil, err
	}

	var available []task.Task
	for _, t := range tasks {
		if t.Assignee == "" {
			available = append(available, t)
		}
	}

	if len(available) == 0 {
		return nil, nil
	}

	sort.Slice(available, func(i, j int) bool {
		ri := cfg.PriorityRank(available[i].Priority)
		rj := cfg.PriorityRank(available[j].Priority)
		if ri != rj {
			return ri < rj
		}
		return available[i].ID < available[j].ID
	})

	picked := available[0]
	return &picked, nil
}
