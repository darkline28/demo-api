package taskservices

import (
	"errors"
	"study/api/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name          string
		input         models.Task
		mockSetUp     func(m *MockTaskRepository, input models.Task)
		mockSetUpUser func(u *MockUserFinder, input models.Task)
		wantErr       bool
		wantErrIs     error
	}{
		{
			name:  "успешное создание задачи",
			input: models.Task{Text: "Test task", Status: "new", UserID: 1},
			mockSetUp: func(m *MockTaskRepository, _ models.Task) {
				m.On("Create", mock.AnythingOfType("*models.Task")).Return(nil)
			},
			mockSetUpUser: func(u *MockUserFinder, input models.Task) {
				u.On("FindByID", input.UserID).
					Return(models.User{ID: input.UserID}, nil)
			},
			wantErr:   false,
			wantErrIs: nil,
		},
		{
			name:  "ошибка — пользователь не найден",
			input: models.Task{Text: "bad task", Status: "new", UserID: 1},
			mockSetUpUser: func(u *MockUserFinder, input models.Task) {
				u.On("FindByID", input.UserID).
					Return(models.User{}, gorm.ErrRecordNotFound)
			},
			wantErr:   true,
			wantErrIs: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			if tt.mockSetUp != nil {
				tt.mockSetUp(mockRepo, tt.input)
			}

			mockUserFinder := new(MockUserFinder)
			if tt.mockSetUpUser != nil {
				tt.mockSetUpUser(mockUserFinder, tt.input)
			}

			service := NewTaskService(mockRepo, mockUserFinder)
			result, err := service.Create(tt.input)

			if tt.wantErr {
				assert.Error(t, err)

				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}

				if tt.wantErrIs == ErrUserNotFound {
					mockRepo.AssertNotCalled(t, "Create", mock.Anything)
				}

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.input.Text, result.Text)
			assert.Equal(t, "new", result.Status)

			mockRepo.AssertExpectations(t)
			mockUserFinder.AssertExpectations(t)
		})
	}
}
func TestListTasks(t *testing.T) {
	tests := []struct {
		name      string
		mockSetUp func(m *MockTaskRepository)
		want      []models.Task
		wantErr   bool
	}{
		{
			name: "успешное создание списка задач",
			mockSetUp: func(m *MockTaskRepository) {
				m.On("FindAll").Return([]models.Task{
					{ID: 1, Text: "test1", Status: "new"},
					{ID: 2, Text: "test2", Status: "done"},
				},
					nil,
				)
			},
			want: []models.Task{
				{ID: 1, Text: "test1", Status: "new"},
				{ID: 2, Text: "test2", Status: "done"},
			},
			wantErr: false,
		},
		{
			name: "ошибка при создании списка задач",
			mockSetUp: func(m *MockTaskRepository) {
				m.On("FindAll").Return([]models.Task{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetUp(mockRepo)
			mockUserFinder := new(MockUserFinder)
			service := NewTaskService(mockRepo, mockUserFinder)
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
		want      models.Task
		wantErr   bool
	}{
		{
			name:  "успешный вывод задачи по id",
			input: 1,
			mockSetUp: func(m *MockTaskRepository, input int) {
				m.On("FindByID", input).Return(models.Task{
					ID: 1, Text: "test1", Status: "new",
				},
					nil,
				)
			},
			want: models.Task{
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
				m.On("FindByID", input).Return(models.Task{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetUp(mockRepo, tt.input)
			mockUserFinder := new(MockUserFinder)
			service := NewTaskService(mockRepo, mockUserFinder)

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
		want        models.Task
		wantErr     bool
	}{
		{
			name:        "успешное обновление задачи",
			inputID:     1,
			inputText:   "new test1",
			inputStatus: "new",
			mockSetUp: func(m *MockTaskRepository) {
				m.On("FindByID", 1).Return(models.Task{
					ID:     1,
					Text:   "old test1",
					Status: "done",
				}, nil)
				m.On("Update", mock.AnythingOfType("*models.Task")).Return(nil)
			},
			want: models.Task{
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
				m.On("FindByID", 2).Return(models.Task{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetUp(mockRepo)
			mockUserFinder := new(MockUserFinder)
			service := NewTaskService(mockRepo, mockUserFinder)
			result, err := service.Update(tt.inputID, tt.inputText, tt.inputStatus)
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
				m.On("FindByID", 1).Return(models.Task{
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
				m.On("FindByID", 2).Return(models.Task{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetUp(mockRepo)
			mockUserFinder := new(MockUserFinder)
			service := NewTaskService(mockRepo, mockUserFinder)
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
