package models

import "time"

type User struct {
	ID        int `gorm:"primaryKey"`
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Tasks     []Task
}

type Task struct {
	ID     int `gorm:"primaryKey"`
	Text   string
	Status string
	UserID int
}
