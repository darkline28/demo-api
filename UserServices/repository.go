package userservices

import (
	"gorm.io/gorm"
)

// UserRepository defines persistence operations for working with users
type UserRepository interface {
	Create(user *User) error
	FindAll() ([]User, error)
	FindByID(id int) (User, error)
	FindByEmail(email string) (User, error)
	Update(user *User) error
	Delete(id int) error
}

type userRepo struct {
	db *gorm.DB
}

// NewUserRepository creates a UserRepository backed by the given gorm.DB instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

// Create implements UserRepository.
func (u *userRepo) Create(user *User) error {
	return u.db.Create(user).Error
}

// FindAll implements UserRepository.
func (u *userRepo) FindAll() ([]User, error) {
	var users []User
	err := u.db.Find(&users).Error
	return users, err
}

// FindByID implements UserRepository.
func (u *userRepo) FindByID(id int) (User, error) {
	var user User
	err := u.db.First(&user, "id = ?", id).Error
	return user, err
}

// FindByEmail implements UserRepository.
func (u *userRepo) FindByEmail(email string) (User, error) {
	var user User
	err := u.db.First(&user, "email = ?", email).Error
	return user, err
}

// Update implements UserRepository.
func (u *userRepo) Update(user *User) error {
	return u.db.Save(user).Error
}

// Delete implements UserRepository.
func (u *userRepo) Delete(id int) error {
	return u.db.Delete(&User{}, "id = ?", id).Error
}
