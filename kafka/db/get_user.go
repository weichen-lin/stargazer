package database

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetUserEmail(driver neo4j.DriverWithContext, name string) (string, error) {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	email, _ := session.ExecuteRead(context.Background(), func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(context.Background(), `
			MATCH (u:User {name: $name})
			RETURN u.email AS email
            `,
			map[string]interface{}{
				"name":    name,
			})

		if err != nil {
			return "", err
		}

		if result.Err() != nil {
			return "", result.Err()
		}

		if result.Next(context.Background()) {
			record := result.Record()
			email, ok := record.Get("email")
			if !ok {
				return 0, fmt.Errorf("error at getting email from record: %v", record)
			}

			return email, nil
		}

		return "", result.Err()
	})

	if email, ok := email.(string); ok {
		return email, nil
	} else {
		return "", fmt.Errorf("error at converting email to striing: %v", email)
	}
}

func GetUserGithubToken(driver neo4j.DriverWithContext, userId string) (string, error) {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	access_token, _ := session.ExecuteRead(context.Background(), func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(context.Background(), `
			MATCH (a:Account {providerAccountId: $userId})
			RETURN a.access_token AS token
            `,
			map[string]interface{}{
				"userId":    userId,
			})

		if err != nil {
			return "", err
		}

		if result.Err() != nil {
			return "", result.Err()
		}

		if result.Next(context.Background()) {
			record := result.Record()
			access_token, ok := record.Get("token")
			if !ok {
				return "", fmt.Errorf("error at getting token from record: %v", record)
			}

			return access_token, nil
		}

		return "", result.Err()
	})

	if access_token, ok := access_token.(string); ok {
		return access_token, nil
	} else {
		return "", fmt.Errorf("error at converting token to striing: %v", access_token)
	}
}