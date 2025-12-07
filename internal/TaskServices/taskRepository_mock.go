package taskservices

import "github.com/stretchr/testify/mock"

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task *Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) FindAll() ([]Task, error) {
	args := m.Called()
	var tasks []Task
	if res := args.Get(0); res != nil {
		tasks = res.([]Task)
	}
	return tasks, args.Error(1)
}

func (m *MockTaskRepository) FindByID(id int) (Task, error) {
	args := m.Called(id)
	task, _ := args.Get(0).(Task)
	return task, args.Error(1)
}

func (m *MockTaskRepository) Update(task *Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
