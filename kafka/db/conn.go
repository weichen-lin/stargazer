package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(src string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  src,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	return db, err
}
