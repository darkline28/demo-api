// Package db contains helpers for working with the application database
package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB opens a PostgreSQL connection and returns a gorm.DB instance
func InitDB() (*gorm.DB, error) {

	dsn := "host=localhost user=project2_user password=pass2 dbname=project2_db port=5433 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil

}
