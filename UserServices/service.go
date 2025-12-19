package userservices

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// UserService describes the business logic of user management
type UserService interface {
	Create(user User) (User, error)
	List() ([]User, error)
	GetByID(id int) (User, error)
	GetByEmail(email string) (User, error)
	Update(id int, email, password string) (User, error)
	Delete(id int) error
}

type userService struct {
	repo UserRepository
}

// NewUserService constructs a UserService using the provided repository
func NewUserService(r UserRepository) UserService {
	return &userService{
		repo: r,
	}
}

// Create implements UserService.
func (u *userService) Create(user User) (User, error) {
	if user.Email == "" {
		return User{}, errors.New("email is empty")
	}
	if user.Password == "" {
		return User{}, errors.New("password is empty")
	}

	_, err := u.repo.FindByEmail(user.Email)
	if err == nil {
		return User{}, errors.New("this email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, err
	}
	if err := u.repo.Create(&user); err != nil {
		return User{}, fmt.Errorf("server error: %w", err)
	}

	return user, nil

}

// List implements UserService.
func (u *userService) List() ([]User, error) {
	users, err := u.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("list users error: %w", err)
	}
	return users, nil
}

// ErrUserNotFound var from Error "user not found"
var ErrUserNotFound = errors.New("user not found")

// GetByID implements UserService.
func (u *userService) GetByID(id int) (User, error) {
	if id <= 0 {
		return User{}, fmt.Errorf("invalid id")
	}
	user, err := u.repo.FindByID(id)
	if err == nil {
		return user, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, ErrUserNotFound
	}

	return User{}, fmt.Errorf("get user by id: %w", err)

}

// GetByEmail implements UserService.
func (u *userService) GetByEmail(email string) (User, error) {
	if email == "" {
		return User{}, fmt.Errorf("invalid email")
	}
	user, err := u.repo.FindByEmail(email)
	if err == nil {
		return user, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, ErrUserNotFound
	}
	return User{}, fmt.Errorf("get user by email: %w", err)
}

// Update implements UserService.
func (u *userService) Update(id int, email string, password string) (User, error) {
	if id <= 0 {
		return User{}, fmt.Errorf("invalid id")
	}

	user, err := u.repo.FindByID(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("server error: %w", err)
	}
	if email == "" && password == "" {
		return User{}, fmt.Errorf("nothing to update")
	}
	if email != "" {
		user.Email = email
	}
	if password != "" {
		user.Password = password
	}

	err = u.repo.Update(&user)
	if err != nil {
		return User{}, fmt.Errorf("server error: %w", err)
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
