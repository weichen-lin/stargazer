package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=johndoe password=randompassword dbname=mydb port=5432 sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	return db, err
}
