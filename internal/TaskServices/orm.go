// Package taskservices defines task domain models used by the service and repository layers
package taskservices

// Task represents a to-do item stored in the database
type Task struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Text   string `json:"text"`
	Status string `json:"status"`
}

// TaskRequest describes the payload used to create a new task via the API
type TaskRequest struct {
	Task string `json:"text"`
}
