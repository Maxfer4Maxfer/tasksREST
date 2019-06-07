package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Storage declares methods that the real storage object should implement
type Storage interface {
	SaveTask(t Task) error
	Task(ID uuid.UUID) (Task, error)
	AllTasks() ([]Task, error)
	DeleteAllTasks() error
}

// Task represents a business object
type Task struct {
	UUID      uuid.UUID `json:"-" gorm:"primary_key"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// Service is a central component of the system. It contains all business logic.
type Service struct {
	storage Storage
}

var (
	// ErrTaskNotFound show that requested task cannot be found
	ErrTaskNotFound = errors.New("Task not found")

	// ErrStorageUnavailable arised if something happened with a storage subsystem
	ErrStorageUnavailable = errors.New("Storage unavailable")
)

// New retunrs new Service
func New(store Storage) *Service {
	s := &Service{
		storage: store,
	}
	return s
}

// NewTask creates a new task and save it in a storage
func (s *Service) NewTask() (string, error) {
	t := Task{
		UUID:      uuid.New(),
		Status:    "created",
		Timestamp: time.Now(),
	}

	err := s.SaveTask(t)
	if err != nil {
		return "", err
	}

	// start gorutine
	go s.taskLifecycle(t.UUID)

	return t.UUID.String(), nil
}

// SaveTask saves a task
func (s *Service) SaveTask(t Task) error {

	err := s.storage.SaveTask(t)
	if err != nil {
		return err
	}
	return nil
}

// Task returns a task with provided guid
func (s *Service) Task(ID string) (Task, error) {
	UUID, err := uuid.Parse(ID)
	if err != nil {
		return Task{}, err
	}

	return s.storage.Task(UUID)
}

// AllTasks returns a task with provided guid
func (s *Service) AllTasks() ([]Task, error) {
	return s.storage.AllTasks()
}

// DeleteAllTasks delete all task in a storage
func (s *Service) DeleteAllTasks() error {
	return s.storage.DeleteAllTasks()
}

func (s *Service) taskLifecycle(UUID uuid.UUID) {
	timer := time.NewTimer(2 * time.Minute)
	// change a status to running
	t := Task{
		UUID:      UUID,
		Status:    "running",
		Timestamp: time.Now(),
	}
	for {
		err := s.SaveTask(t)
		if err == nil {
			break
		}
	}

	// after 2 minutes mark the task by the "finished" status
	<-timer.C

	t = Task{
		UUID:      UUID,
		Status:    "finished",
		Timestamp: time.Now(),
	}
	for {
		err := s.SaveTask(t)
		if err == nil {
			break
		}
	}
}
