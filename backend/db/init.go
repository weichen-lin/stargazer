package db

import (
	"context"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Database struct {
	Driver  neo4j.DriverWithContext
	Timeout int
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

	err = InitFullTextIndex(driver)
	if err != nil {
		panic(err)
	}

	err = InitUserIndex(driver)
	if err != nil {
		panic(err)
	}

	err = InitRepoLanguagesIndex(driver)
	if err != nil {
		panic(err)
	}

	err = InitRepoIdIndex(driver)
	if err != nil {
		panic(err)
	}


	return &Database{
		Driver:  driver,
		Timeout: 5,
	}
}

func InitFullTextIndex(driver neo4j.DriverWithContext) error {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		_, err := transaction.Run(context.Background(),
			"CREATE FULLTEXT INDEX REPOSITORY_FULL_TEXT_SEARCH IF NOT EXISTS "+
				"FOR (r:Repository) ON EACH [r.repo_name, r.description]",
			map[string]interface{}{},
		)

		return nil, err
	})

	if err != nil {
		return err
	}

	return nil
}

func InitUserIndex(driver neo4j.DriverWithContext) error {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		_, err := transaction.Run(context.Background(),
			"CREATE INDEX USER_EMAIL IF NOT EXISTS "+
				"FOR (u:User) ON (u.email);",
			map[string]interface{}{},
		)

		return nil, err
	})

	if err != nil {
		return err
	}

	return nil
}

func InitRepoLanguagesIndex(driver neo4j.DriverWithContext) error {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		_, err := transaction.Run(context.Background(),
			"CREATE INDEX REPOSITORY_LANGUAGE IF NOT EXISTS "+
				"FOR (r:Reposiotry) ON (r.language);",
			map[string]interface{}{},
		)

		return nil, err
	})

	if err != nil {
		return err
	}

	return nil
}

func InitRepoIdIndex(driver neo4j.DriverWithContext) error {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		_, err := transaction.Run(context.Background(),
			"CREATE INDEX REPOSITORY_ID IF NOT EXISTS "+
				"FOR (r:Reposiotry) ON (r.repo_id);",
			map[string]interface{}{},
		)

		return nil, err
	})

	if err != nil {
		return err
	}

	return nil
}

