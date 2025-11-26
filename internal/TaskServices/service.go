package taskservices

import (
	"errors"
	"fmt"
)

// TaskService describes business logic for managing tasks
type TaskService interface {
	Create(text Task) (Task, error)
	List() ([]Task, error)
	GetByID(id int) (Task, error)
	Update(id int, text, status string) (Task, error)
	Delete(id int) error
}

type taskService struct {
	repo TaskRepository
}

// NewTaskService constructs a TaskService using the provided repository
func NewTaskService(r TaskRepository) TaskService {
	return &taskService{repo: r}
}

func (s *taskService) Create(text Task) (Task, error) {
	if text.Text == "" {
		return Task{}, errors.New("text required")
	}
	task := Task{
		Text:   text.Text,
		Status: "new",
	}
	if err := s.repo.Create(&task); err != nil {
		return Task{}, fmt.Errorf("server error: %w", err)
	}

	return task, nil

}

// List implements TaskService.
func (s *taskService) List() ([]Task, error) {
	tasks, err := s.repo.FindAll()
	if err != nil {
		return []Task{}, fmt.Errorf("server error: %w", err)
	}
	return tasks, nil
}

// GetByID implements TaskService.
func (s *taskService) GetByID(id int) (Task, error) {
	if id <= 0 {
		return Task{}, errors.New("invalid id")
	}
	task, err := s.repo.FindByID(id)
	if err != nil {
		return Task{}, fmt.Errorf("server error: %w", err)
	}
	return task, nil
}

// Update implements TaskService.
func (s *taskService) Update(id int, text string, status string) (Task, error) {
	if id <= 0 {
		return Task{}, errors.New("invalid id")
	}
	task, err := s.repo.FindByID(id)
	if err != nil {
		return Task{}, fmt.Errorf("server error: %w", err)
	}
	if text != "" {
		task.Text = text
	}
	if status != "" {
		if !isValidationStatus(status) {
			return Task{}, errors.New("invalid status")
		}
		task.Status = status
	}
	if text == "" && status == "" {
		return Task{}, errors.New("nothing to update")
	}
	err = s.repo.Update(&task)
	if err != nil {
		return Task{}, fmt.Errorf("server error: %w", err)
	}
	return task, nil
}

// Delete implements TaskService.
func (s *taskService) Delete(id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	_, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	if err = s.repo.Delete(id); err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

func isValidationStatus(s string) bool {
	if s == "new" || s == "in progress" || s == "done" {
		return true
	}
	return false
}
