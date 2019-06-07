package gorm

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"github.com/jinzhu/gorm"

	"tasksREST/pkg/service"
)

// // Task represents a business object
// type Task struct {
// 	UUID      uuid.UUID `gorm:"primary_key"`
// 	Status    string
// 	Timestamp time.Time
// }

// Storage implements mySQL storage for nodes
type Storage struct {
	DB  *gorm.DB
	DSN string
}

// New create in memory repository for storing nodes
func New(dsn string) (*Storage, func()) {

	s := &Storage{
		DB:  nil,
		DSN: dsn,
	}

	ctx, cancel := context.WithCancel(context.Background())
	go s.connectToDB(ctx)

	return s, cancel
}

func (s *Storage) connectToDB(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			if s.DB == nil {
				db, err := gorm.Open("mysql", s.DSN)
				if err != nil {
					fmt.Println("db", "mysql", "message", "got an error", "err", err)
				} else {
					fmt.Println("db", "mysql", "message", "connection is established")
					s.DB = db
					s.DB.AutoMigrate(&service.Task{})
				}

			}
			if s.DB != nil {
				err := s.DB.DB().Ping()
				if err != nil {
					fmt.Println("db", "mysql", "message", "lost connection to the db", "err", err)
				}
			}
		case <-ctx.Done():
			ticker.Stop()
			s.DB.Close()
		}
	}
}

// SaveTask creates or saves a task in a database
func (s *Storage) SaveTask(t service.Task) error {
	if s.DB != nil {
		s.DB.Save(&t)
		return nil
	} else {
		return service.ErrStorageUnavailable
	}
}

// Task returns a task by UUID
func (s *Storage) Task(ID uuid.UUID) (service.Task, error) {
	t := service.Task{}
	if s.DB != nil {
		var count int
		s.DB.Where("UUID = ?", ID).First(&t).Count(&count)
		if count == 0 {
			return t, service.ErrTaskNotFound
		}
		return t, nil
	} else {
		return t, service.ErrStorageUnavailable
	}
}

// AllTasks returns all tasks
func (s *Storage) AllTasks() ([]service.Task, error) {
	tasks := []service.Task{}

	if s.DB != nil {
		s.DB.Find(&tasks)
		return tasks, nil
	} else {
		return tasks, service.ErrStorageUnavailable
	}
}

// DeleteAllTasks deletes all tasks
func (s *Storage) DeleteAllTasks() error {
	s.DB.Where("UUID LIKE ?", "%").Delete(service.Task{})
	return nil
}
