package neo4j_kafka

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func InitializeConstraints(driver neo4j.DriverWithContext) error {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	constraint := "CREATE CONSTRAINT FOR (u:User) REQUIRE u.user_id IS UNIQUE;"
	_, err := session.ExecuteWrite(context.Background(), func(transaction neo4j.ManagedTransaction) (any, error) {
		_, err := transaction.Run(context.Background(), constraint, nil)
		return nil, err
	})

	if neoErr, ok := err.(*neo4j.Neo4jError); ok {
		switch neoErr.Code {
		case "Neo.ClientError.Schema.EquivalentSchemaRuleAlreadyExists":
			return nil
		default:
			return err
		}
	} else {
		return err
	}
}
