package taskservices

import (
	"study/api/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockTaskRepository from test
type MockTaskRepository struct {
	mock.Mock
}
type MockUserFinder struct {
	mock.Mock
}

// FindByUserID mock method from test
func (m *MockUserFinder) FindByID(id int) (models.User, error) {
	args := m.Called(id)
	user, _ := args.Get(0).(models.User)
	return user, args.Error(1)
}

// Create mock method from test
func (m *MockTaskRepository) Create(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

// FindAll mock method from test
func (m *MockTaskRepository) FindAll() ([]models.Task, error) {
	args := m.Called()
	var tasks []models.Task
	if res := args.Get(0); res != nil {
		tasks = res.([]models.Task)
	}
	return tasks, args.Error(1)
}

// FindByID mock method from test
func (m *MockTaskRepository) FindByID(id int) (models.Task, error) {
	args := m.Called(id)
	task, _ := args.Get(0).(models.Task)
	return task, args.Error(1)
}

// Update mock method from test
func (m *MockTaskRepository) Update(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

// Delete mock method from test
func (m *MockTaskRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTaskRepository) FindByUserID(userID int) ([]models.Task, error) {
	args := m.Called(userID)
	tasks, _ := args.Get(0).([]models.Task)
	return tasks, args.Error(1)
}
