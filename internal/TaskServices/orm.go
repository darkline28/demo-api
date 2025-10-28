package taskservices

type Task struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Text   string `json:"text"`
	Status string `json:"status"`
}
type TaskRequest struct {
	Task string `json:"text"`
}
