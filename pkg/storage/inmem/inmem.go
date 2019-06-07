package inmem

import (
	"github.com/google/uuid"
	"tasksREST/pkg/service"
)

// Storage stores objects in memory
type Storage struct {
	tasks map[uuid.UUID]service.Task
}

// New returns a storage object
func New(dsn string) (*Storage, func()) {
	s := &Storage{
		tasks: make(map[uuid.UUID]service.Task),
	}
	return s, func() {}
}

// SaveTask saves a task
func (s *Storage) SaveTask(t service.Task) error {
	s.tasks[t.UUID] = t
	return nil
}

// Task returns a task from the storage
func (s *Storage) Task(ID uuid.UUID) (service.Task, error) {
	_, ok := s.tasks[ID]
	if !ok {
		return service.Task{}, service.ErrTaskNotFound
	}
	return s.tasks[ID], nil
}

// AllTasks returns all task
func (s *Storage) AllTasks() ([]service.Task, error) {
	r := []service.Task{}
	for _, t := range s.tasks {
		r = append(r, t)
	}
	return r, nil
}

// DeleteAllTasks deletes all tasks
func (s *Storage) DeleteAllTasks() error {
	s.tasks = make(map[uuid.UUID]service.Task)
	return nil
}
