package userservices

import (
	"errors"
	"fmt"
	"study/api/internal/models"

	"gorm.io/gorm"
)

// UserService describes the business logic of user management
type UserService interface {
	Create(user models.User) (models.User, error)
	List() ([]models.User, error)
	GetByID(id int) (models.User, error)
	GetByEmail(email string) (models.User, error)
	Update(id int, email, password string) (models.User, error)
	Delete(id int) error
	GetTasksForUser(userID int) ([]models.Task, error)
}

type TaskFinder interface {
	FindByUserID(userID int) ([]models.Task, error)
}

type userService struct {
	repo     UserRepository
	taskRepo TaskFinder
}

// NewUserService constructs a UserService using the provided repository
func NewUserService(r UserRepository, t TaskFinder) UserService {
	return &userService{
		repo:     r,
		taskRepo: t,
	}
}

// GetTasksForUser implements UserService.
func (u *userService) GetTasksForUser(userID int) ([]models.Task, error) {
	if userID <= 0 {
		return []models.Task{}, fmt.Errorf("invalid id")
	}
	_, err := u.repo.FindByID(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []models.Task{}, ErrUserNotFound
	}
	if err != nil {
		return []models.Task{}, fmt.Errorf("server error: %w", err)
	}

	tasks, tasksErr := u.taskRepo.FindByUserID(userID)
	if tasksErr != nil {
		return []models.Task{}, fmt.Errorf("get tasks for user: %w", tasksErr)
	}

	return tasks, nil
}

// Create implements UserService.
func (u *userService) Create(user models.User) (models.User, error) {
	if user.Email == "" {
		return models.User{}, errors.New("email is empty")
	}
	if user.Password == "" {
		return models.User{}, errors.New("password is empty")
	}

	_, err := u.repo.FindByEmail(user.Email)
	if err == nil {
		return models.User{}, errors.New("this email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, err
	}
	if err := u.repo.Create(&user); err != nil {
		return models.User{}, fmt.Errorf("server error: %w", err)
	}

	return user, nil

}

// List implements UserService.
func (u *userService) List() ([]models.User, error) {
	users, err := u.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("list users error: %w", err)
	}
	return users, nil
}

// ErrUserNotFound var from Error "user not found"
var ErrUserNotFound = errors.New("user not found")

// GetByID implements UserService.
func (u *userService) GetByID(id int) (models.User, error) {
	if id <= 0 {
		return models.User{}, fmt.Errorf("invalid id")
	}
	user, err := u.repo.FindByID(id)
	if err == nil {
		return user, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, ErrUserNotFound
	}

	return models.User{}, fmt.Errorf("get user by id: %w", err)

}

// GetByEmail implements UserService.
func (u *userService) GetByEmail(email string) (models.User, error) {
	if email == "" {
		return models.User{}, fmt.Errorf("invalid email")
	}
	user, err := u.repo.FindByEmail(email)
	if err == nil {
		return user, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, ErrUserNotFound
	}
	return models.User{}, fmt.Errorf("get user by email: %w", err)
}

// Update implements UserService.
func (u *userService) Update(id int, email string, password string) (models.User, error) {
	if id <= 0 {
		return models.User{}, fmt.Errorf("invalid id")
	}

	user, err := u.repo.FindByID(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("server error: %w", err)
	}
	if email == "" && password == "" {
		return models.User{}, fmt.Errorf("nothing to update")
	}
	if email != "" {
		user.Email = email
	}
	if password != "" {
		user.Password = password
	}

	err = u.repo.Update(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("server error: %w", err)
	}
	return user, nil
}

// Delete implements UserService.
func (u *userService) Delete(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid id")
	}
	_, err := u.repo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	if err = u.repo.Delete(id); err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}
