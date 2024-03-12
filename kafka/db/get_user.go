package database

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/workflow"
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
				"name": name,
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
		return "", fmt.Errorf("error at converting email to striing")
	}
}

func GetUserGithubToken(driver neo4j.DriverWithContext, userId string) (string, error) {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	access_token, err := session.ExecuteRead(context.Background(), func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(context.Background(), `
			MATCH (a:Account {providerAccountId: $userId})
			RETURN a.access_token AS token
            `,
			map[string]interface{}{
				"userId": userId,
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

	if err != nil {
		fmt.Println("Error getting token from neo4j:", err)
		return "", err
	}

	if access_token, ok := access_token.(string); ok {
		return access_token, nil
	} else {
		return "", fmt.Errorf("error at converting token to string")
	}
}

func ConfirmVectorize(driver neo4j.DriverWithContext, info *workflow.SyncUserStarMsg) error {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(context.Background(), `
			MATCH (u:User { name: $name })-[s:STARS]->(r:Repository { repo_id: $repo_id })
			SET s.is_vectorized = $isVectorized
			RETURN s{.*}
            `,
			map[string]interface{}{
				"name":         info.UserName,
				"repo_id":      info.RepoId,
				"isVectorized": true,
			})

		if err != nil {
			return "", err
		}

		if result.Err() != nil {
			return "", result.Err()
		}

		return "", result.Err()
	})

	if err != nil {
		fmt.Println("Error make vectorize success from neo4j:", err)
		return err
	}

	return nil
}
