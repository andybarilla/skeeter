package store

import (
	"github.com/andybarilla/skeeter/internal/task"
)

type Filter struct {
	Status   string
	Priority string
	Assignee string
	Tag      string
}

type Store interface {
	Init(projectName string) error
	List(filter Filter) ([]task.Task, error)
	Get(id string) (*task.Task, error)
	Create(t *task.Task) error
	Update(t *task.Task) error
	NextID() (string, error)
}
