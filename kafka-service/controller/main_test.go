package controller

import (
	"os"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/db"
)

var testDB *db.Database
var testController *Controller

func NewTestDatabase() *db.Database {
	driver, err := neo4j.NewDriverWithContext(
		"neo4j://localhost:7687",
		neo4j.BasicAuth("neo4j", "password", ""),
	)

	if err != nil {
		panic(err)
	}

	return &db.Database{
		Driver: driver,
	}
}

func TestMain(m *testing.M) {
	testDB = NewTestDatabase()

	testController = &Controller{
		db: testDB,
	}

	os.Exit(m.Run())
}
