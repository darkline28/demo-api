package taskservices

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name      string
		input     Task
		mockSetUp func(m *MockTaskRepository, input Task)
		wantErr   bool
	}{
		{
			name:  "успешное создание задачи",
			input: Task{Text: "Test task", Status: "new"},
			mockSetUp: func(m *MockTaskRepository, input Task) {
				m.On("Create", mock.AnythingOfType("*taskservices.Task")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "ошибка при создании",
			input: Task{Text: "bad task", Status: "new"},
			mockSetUp: func(m *MockTaskRepository, input Task) {
				m.On("Create", mock.AnythingOfType("*taskservices.Task")).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetUp(mockRepo, tt.input)

			service := NewTaskService(mockRepo)
			result, err := service.Create(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Text, result.Text)
				assert.Equal(t, "new", result.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestListTasks(t *testing.T) {
	tests := []struct {
		name      string
		mockSetUp func(m *MockTaskRepository)
		want      []Task
		wantErr   bool
	}{
		{
			name: "успешное создание списка задач",
			mockSetUp: func(m *MockTaskRepository) {
				m.On("FindAll").Return([]Task{
					{ID: 1, Text: "test1", Status: "new"},
					{ID: 2, Text: "test2", Status: "done"},
				},
					nil,
				)
			},
			want: []Task{
				{ID: 1, Text: "test1", Status: "new"},
				{ID: 2, Text: "test2", Status: "done"},
			},
			wantErr: false,
		},
		{
			name: "ошибка при создании списка задач",
			mockSetUp: func(m *MockTaskRepository) {
				m.On("FindAll").Return([]Task{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetUp(mockRepo)

			service := NewTaskService(mockRepo)
			result, err := service.List()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetByIDTasks(t *testing.T) {
	tests := []struct {
		name      string
		input     int
		mockSetUp func(m *MockTaskRepository, input int)
		want      Task
		wantErr   bool
	}{
		{
			name:  "успешный вывод задачи по id",
			input: 1,
			mockSetUp: func(m *MockTaskRepository, input int) {
				m.On("FindByID", input).Return(Task{
					ID: 1, Text: "test1", Status: "new",
				},
					nil,
				)
			},
			want: Task{
				ID:     1,
				Text:   "test1",
				Status: "new",
			},
			wantErr: false,
		},
		{
			name:  "ошибка вывода задачи по id",
			input: 2,
			mockSetUp: func(m *MockTaskRepository, input int) {
				m.On("FindByID", input).Return(Task{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetUp(mockRepo, tt.input)

			service := NewTaskService(mockRepo)
			result, err := service.GetByID(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
			mockRepo.AssertExpectations(t)
		})

	}
}

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name        string
		inputID     int
		inputText   string
		inputStatus string
		mockSetUp   func(m *MockTaskRepository)
		want        Task
		wantErr     bool
	}{
		{
			name:        "успешное обновление задачи",
			inputID:     1,
			inputText:   "new test1",
			inputStatus: "new",
			mockSetUp: func(m *MockTaskRepository) {
				m.On("FindByID", 1).Return(Task{
					ID:     1,
					Text:   "old test1",
					Status: "done",
				}, nil)
				m.On("Update", mock.AnythingOfType("*taskservices.Task")).Return(nil)
			},
			want: Task{
				ID:     1,
				Text:   "new test1",
				Status: "new",
			},
			wantErr: false,
		},
		{
			name:        "ошибка при обновлении задачи",
			inputID:     2,
			inputText:   "old test2",
			inputStatus: "done",
			mockSetUp: func(m *MockTaskRepository) {
				m.On("FindByID", 2).Return(Task{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockrepo := new(MockTaskRepository)
			tt.mockSetUp(mockrepo)

			service := NewTaskService(mockrepo)
			result, err := service.Update(tt.inputID, tt.inputText, tt.inputStatus)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
			mockrepo.AssertExpectations(t)
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name      string
		inputID   int
		mockSetUp func(m *MockTaskRepository)
		wantErr   bool
	}{
		{
			name:    "успешное удаление задачи",
			inputID: 1,
			mockSetUp: func(m *MockTaskRepository) {
				m.On("FindByID", 1).Return(Task{
					ID:     1,
					Text:   "task1",
					Status: "done",
				}, nil)
				m.On("Delete", 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "ошибка при удалении задачи",
			inputID: 2,
			mockSetUp: func(m *MockTaskRepository) {
				m.On("FindByID", 2).Return(Task{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetUp(mockRepo)

			service := NewTaskService(mockRepo)
			err := service.Delete(tt.inputID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}

}
