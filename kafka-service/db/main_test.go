package db

import (
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var db *Database

func NewTestDatabase() *Database {
	driver, err := neo4j.NewDriverWithContext(
		"neo4j://localhost:7687",
		neo4j.BasicAuth("neo4j", "password", ""),
	)

	if err != nil {
		panic(err)
	}

	return &Database{
		driver: driver,
	}
}

func TestMain(m *testing.M) {
	db = NewDatabase()
}