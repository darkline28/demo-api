package taskservices

import (
	"errors"
	"fmt"
	"study/api/internal/models"

	"gorm.io/gorm"
)

// TaskService describes business logic for managing tasks
type TaskService interface {
	Create(text models.Task) (models.Task, error)
	List() ([]models.Task, error)
	GetByID(id int) (models.Task, error)
	Update(id int, text, status string) (models.Task, error)
	Delete(id int) error
}

type UserFinder interface {
	FindByID(id int) (models.User, error)
}

type taskService struct {
	repo     TaskRepository
	userRepo UserFinder
}

// NewTaskService constructs a TaskService using the provided repository
func NewTaskService(r TaskRepository, u UserFinder) TaskService {
	return &taskService{repo: r, userRepo: u}
}

// ErrUserNotFound var from Error "user not found"
var ErrUserNotFound = errors.New("user not found")

func (s *taskService) Create(text models.Task) (models.Task, error) {
	if text.Text == "" {
		return models.Task{}, errors.New("text required")
	}
	if text.UserID <= 0 {
		return models.Task{}, fmt.Errorf("invalid user id")
	}
	task := models.Task{
		Text:   text.Text,
		Status: "new",
		UserID: text.UserID,
	}
	_, err := s.userRepo.FindByID(text.UserID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Task{}, ErrUserNotFound
	}

	if err != nil {
		return models.Task{}, fmt.Errorf("get user by id: %w", err)
	}

	if err := s.repo.Create(&task); err != nil {
		return models.Task{}, fmt.Errorf("server error: %w", err)
	}

	return task, nil

}

// List implements TaskService.
func (s *taskService) List() ([]models.Task, error) {
	tasks, err := s.repo.FindAll()
	if err != nil {
		return []models.Task{}, fmt.Errorf("server error: %w", err)
	}
	return tasks, nil
}

// GetByID implements TaskService.
func (s *taskService) GetByID(id int) (models.Task, error) {
	if id <= 0 {
		return models.Task{}, errors.New("invalid id")
	}
	task, err := s.repo.FindByID(id)
	if err != nil {
		return models.Task{}, fmt.Errorf("server error: %w", err)
	}
	return task, nil
}

// Update implements TaskService.
func (s *taskService) Update(id int, text string, status string) (models.Task, error) {
	if id <= 0 {
		return models.Task{}, errors.New("invalid id")
	}
	task, err := s.repo.FindByID(id)
	if err != nil {
		return models.Task{}, fmt.Errorf("server error: %w", err)
	}
	if text != "" {
		task.Text = text
	}
	if status != "" {
		if !isValidationStatus(status) {
			return models.Task{}, errors.New("invalid status")
		}
		task.Status = status
	}
	if text == "" && status == "" {
		return models.Task{}, errors.New("nothing to update")
	}
	err = s.repo.Update(&task)
	if err != nil {
		return models.Task{}, fmt.Errorf("server error: %w", err)
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
