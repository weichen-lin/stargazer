package neo4j_kafka

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func handleNeo4jError(err error) error {
	if err == nil {
		return nil
	}

	if neoErr, ok := err.(*neo4j.Neo4jError); ok {
		switch neoErr.Code {
		case "Neo.ClientError.Schema.ConstraintValidationFailed":
			return errors.New("user already exists")
		case "Neo.ClientError.Schema.EquivalentSchemaRuleAlreadyExists":
			return nil
		default:
			return err
		}
	} else {
		return err
	}
}
