package neo4j_kafka

import (
	"context"
	"errors"

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

	if neoErr, ok := err.(*neo4j.Neo4jError); ok {
		switch neoErr.Code {
		case "Neo.ClientError.Schema.ConstraintValidationFailed":
			return errors.New("User already exists")
		default:
			return err
		}
	} else {
		return err
	}

}
