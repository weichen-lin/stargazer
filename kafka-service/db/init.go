package db

import (
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Database struct {
	driver neo4j.DriverWithContext
}

func NewDatabase() *Database {
	neo4j_url := os.Getenv("NEO4J_URL")
	neo4j_password := os.Getenv("NEO4J_PASSWORD")

	driver, err := neo4j.NewDriverWithContext(
		neo4j_url,
		neo4j.BasicAuth("neo4j", neo4j_password, ""),
	)

	if err != nil {
		panic(err)
	}

	return &Database{
		driver: driver,
	}
}
