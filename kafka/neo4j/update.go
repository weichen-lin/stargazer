package neo4jOpeartion

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/workflow"
)

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
