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

func GetUserNotVectorize(driver neo4j.DriverWithContext, userName string) ([]int64, error) {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	records, err := session.ExecuteRead(context.Background(), func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(context.Background(), `
				MATCH (u:User {name: $name})-[s:STARS]-(r:Repository)
				WHERE s.is_vectorized = FALSE or s.is_vectorized IS NULL
				RETURN r.repo_id as repo_id
            `,
			map[string]interface{}{
				"name": userName,
			})

		if err != nil {
			return nil, err
		}

		if result.Err() != nil {
			return nil, result.Err()
		}

		collects, err := result.Collect(context.Background())
		if err != nil {
			return nil, err
		}

		stars := make([]int64, len(collects))

		for i, record := range collects {
			repo_id, ok := record.Get("repo_id")
			if !ok {
				return nil, fmt.Errorf("error at getting repo_id from record: %v", record)
			}

			stars[i] = repo_id.(int64)
		}

		return stars, result.Err()
	})

	if err != nil {
		fmt.Println("Error make vectorize success from neo4j:", err)
		return nil, err
	}

	if _, ok := records.([]int64); !ok {
		return nil, fmt.Errorf("error at converting stars to []int")
	}

	return records.([]int64), err
}
