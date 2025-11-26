package taskservices

import (
	"gorm.io/gorm"
)

// TaskRepository defines persistence operations for working with tasks
type TaskRepository interface {
	Create(task *Task) error
	FindAll() ([]Task, error)
	FindByID(id int) (Task, error)
	Update(task *Task) error
	Delete(id int) error
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

func (r *taskRepo) Create(task *Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepo) FindAll() ([]Task, error) {
	var tasks []Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *taskRepo) FindByID(id int) (Task, error) {
	var task Task
	err := r.db.First(&task, "id = ?", id).Error
	return task, err
}

func (r *taskRepo) Update(task *Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepo) Delete(id int) error {
	return r.db.Delete(&Task{}, "id = ?", id).Error
}
