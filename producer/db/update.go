package db

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/workflow"
)

func (db *Database) ConfirmVectorize(info *workflow.SyncUserStar) error {
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

func (db *Database) WriteResultAtCrontab(email, status string) error {
	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	records, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})
			MERGE (u)-[h:HAS_CRONTAB]-(c:Crontab)
			SET c += {
				created_at: COALESCE(c.created_at, datetime()),
				last_trigger_time: datetime(),
				status: $status
			}
			RETURN c.last_trigger_time AS last_trigger_time;
            `, map[string]any{
			"email":  email,
			"status": status,
		})

		if err != nil {
			return nil, err
		}

		records, _ := result.Collect(context.Background())
		if err != nil {
			return nil, err
		}

		return records, nil
	})

	if err != nil {
		return err
	}

	crontabs, ok := records.([]*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting crontabs records to []*neo4j.Record")
	}

	if len(crontabs) == 0 {
		return fmt.Errorf("user not found")
	}

	record := crontabs[0].AsMap()
	_, ok = record["last_trigger_time"].(time.Time)

	if !ok {
		return fmt.Errorf("error at convert last_trigger_time from record: %v", record)
	}

	return nil
}
