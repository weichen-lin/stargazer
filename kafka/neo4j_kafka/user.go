package neo4j_kafka

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/workflow"
)

func CreateUser(driver neo4j.DriverWithContext, user workflow.User) error {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(context.Background(), `
			CREATE (u:User {
			id: apoc.create.uuid(),
			user_id: $user_id,
			name: $name,
			token: $token,
			is_sync: false,
			createdAt: datetime(),
			updatedAt: datetime()
			});
			`,
			map[string]interface{}{
				"user_id": user.ID,
				"name":    user.Login,
				"token":   "",
			})

		if err != nil {
			return nil, err
		}

		return nil, result.Err()
	})

	err = handleNeo4jError(err)
	return err
}
