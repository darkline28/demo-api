// Package userservices defines task domain models used by the service and repository layers
package userservices

import "time"

// User fields stored in the database
type User struct {
	ID        int        `gorm:"primaryKey" json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
