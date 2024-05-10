package db

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/workflow"
)

func (db *database) ConfirmVectorize(info *workflow.SyncUserStar) error {
	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(context.Background(), `
			MATCH (u:User { email: $email })-[s:STARS]->(r:Repository { repo_id: $repo_id })
			SET s.is_vectorized = $isVectorized
			RETURN s{.*}
            `,
			map[string]interface{}{
				"email":        info.Email,
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
