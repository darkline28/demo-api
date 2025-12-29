package taskservices

import (
	"study/api/internal/models"

	"gorm.io/gorm"
)

// TaskRepository defines persistence operations for working with tasks
type TaskRepository interface {
	Create(task *models.Task) error
	FindAll() ([]models.Task, error)
	FindByID(id int) (models.Task, error)
	Update(task *models.Task) error
	Delete(id int) error
	FindByUserID(userID int) ([]models.Task, error)
}

type taskRepo struct {
	db *gorm.DB
}

// NewTaskRepository creates a TaskRepository backed by the given gorm.DB instance
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepo{
		db: db,
	}
}

// FindByUserID implements TaskRepository.
func (r *taskRepo) FindByUserID(userID int) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Find(&tasks, "user_id = ?", userID).Error
	return tasks, err
}

func (r *taskRepo) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepo) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *taskRepo) FindByID(id int) (models.Task, error) {
	var task models.Task
	err := r.db.First(&task, "id = ?", id).Error
	return task, err
}

func (r *taskRepo) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepo) Delete(id int) error {
	return r.db.Delete(&models.Task{}, "id = ?", id).Error
}
