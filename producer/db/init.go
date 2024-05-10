package db

import (
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/workflow"
)

type Database interface {
	// create
	CreateRepository(repo *Repository, email string) error
	UpdateCrontab(hour int, email string) error

	// read
	GetUser(email string) (*User, error)
	GetUserConfig(email string) (*Config, error)
	GetUserNotVectorize(email string) ([]int64, error)
	GetAllUserCrontab() ([]Crontab, error)

	// update
	ConfirmVectorize(info *workflow.SyncUserStar) error
}

type database struct {
	driver neo4j.DriverWithContext
}

func NewDatabase() Database {
	neo4j_url := os.Getenv("NEO4J_URL")
	neo4j_password := os.Getenv("NEO4J_PASSWORD")

	driver, err := neo4j.NewDriverWithContext(
		neo4j_url,
		neo4j.BasicAuth("neo4j", neo4j_password, ""),
	)

	if err != nil {
		panic(err)
	}

	return &database{
		driver: driver,
	}
}
