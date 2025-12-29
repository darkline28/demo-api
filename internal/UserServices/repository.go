package userservices

import (
	"study/api/internal/models"

	"gorm.io/gorm"
)

// UserRepository defines persistence operations for working with users
type UserRepository interface {
	Create(user *models.User) error
	FindAll() ([]models.User, error)
	FindByID(id int) (models.User, error)
	FindByEmail(email string) (models.User, error)
	Update(user *models.User) error
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
func (u *userRepo) Create(user *models.User) error {
	return u.db.Create(user).Error
}

// FindAll implements UserRepository.
func (u *userRepo) FindAll() ([]models.User, error) {
	var users []models.User
	err := u.db.Find(&users).Error
	return users, err
}

// FindByID implements UserRepository.
func (u *userRepo) FindByID(id int) (models.User, error) {
	var user models.User
	err := u.db.First(&user, "id = ?", id).Error
	return user, err
}

// FindByEmail implements UserRepository.
func (u *userRepo) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := u.db.First(&user, "email = ?", email).Error
	return user, err
}

// Update implements UserRepository.
func (u *userRepo) Update(user *models.User) error {
	return u.db.Save(user).Error
}

// Delete implements UserRepository.
func (u *userRepo) Delete(id int) error {
	return u.db.Delete(&models.User{}, "id = ?", id).Error
}
