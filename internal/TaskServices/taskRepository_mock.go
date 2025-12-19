package taskservices

import "github.com/stretchr/testify/mock"

// MockTaskRepository from test
type MockTaskRepository struct {
	mock.Mock
}

// Create mock method from test
func (m *MockTaskRepository) Create(task *Task) error {
	args := m.Called(task)
	return args.Error(0)
}

// FindAll mock method from test
func (m *MockTaskRepository) FindAll() ([]Task, error) {
	args := m.Called()
	var tasks []Task
	if res := args.Get(0); res != nil {
		tasks = res.([]Task)
	}
	return tasks, args.Error(1)
}

// FindByID mock method from test
func (m *MockTaskRepository) FindByID(id int) (Task, error) {
	args := m.Called(id)
	task, _ := args.Get(0).(Task)
	return task, args.Error(1)
}

// Update mock method from test
func (m *MockTaskRepository) Update(task *Task) error {
	args := m.Called(task)
	return args.Error(0)
}

// Delete mock method from test
func (m *MockTaskRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
