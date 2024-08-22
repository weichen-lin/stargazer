package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/domain"
)

var ErrNotFoundCrontab = errors.New("crontab not found")

func (db *Database) GetCrontab(ctx context.Context) (*domain.Crontab, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[h:HAS_CRONTAB]-(c:Crontab)
			RETURN {
				trigger_at: c.trigger_at,
				created_at: c.created_at,
				updated_at: c.updated_at,
				status: c.status,
				last_triggered_at: c.last_triggered_at,
				version: c.version
			} as crontab
			`,
			map[string]interface{}{
				"email": email,
			})

		if err != nil {
			fmt.Println("error at read crontab: ", err)
			return nil, err
		}
		record, err := result.Single(context.Background())
		return record, err
	})

	if err != nil {
		return nil, ErrNotFoundCrontab
	}

	record, ok := result.(*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting crontab records to *neo4j.Record")
	}

	crontabMap := record.AsMap()

	crontab, ok := crontabMap["crontab"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error converting record to map")
	}

	cron, err := domain.FromCrontabEntity(
		&domain.CrontabEntity{
			TriggeredAt:     getString(crontab["triggered_at"]),
			CreatedAt:       getString(crontab["created_at"]),
			UpdatedAt:       getString(crontab["updated_at"]),
			Status:          getString(crontab["status"]),
			LastTriggeredAt: getString(crontab["last_triggered_at"]),
			Version:         getInt64(crontab["version"]),
		},
	)

	if err != nil {
		return nil, err
	}

	return cron, nil
}

func (db *Database) CreateCrontab(ctx context.Context, crontab *domain.Crontab) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	entity := crontab.ToCrontabEntity()

	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	result, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})
			MERGE (u)-[h:HAS_CRONTAB]-(c:Crontab)
			Set c += {
				created_at: $created_at,
				triggered_at: $triggered_at,
				updated_at: $updated_at,
				status: $status,
				last_triggered_at: $last_triggered_at,
				version: $version
			}
			RETURN c.created_at AS created_at;
            `, map[string]any{
			"email":             email,
			"created_at":        entity.CreatedAt,
			"updated_at":        entity.UpdatedAt,
			"triggered_at":      entity.TriggeredAt,
			"last_triggered_at": entity.LastTriggeredAt,
			"status":            entity.Status,
			"version":           entity.Version,
		})

		if err != nil {
			fmt.Println("error at read crontab: ", err)
			return nil, err
		}

		record, err := result.Single(context.Background())
		return record, err
	})

	if err != nil {
		return err
	}

	record, ok := result.(*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting crontab records to []*neo4j.Record")
	}

	resultMap := record.AsMap()

	createdAt, ok := resultMap["created_at"].(string)
	if !ok {
		return fmt.Errorf("error converting record to map")
	}

	if entity.CreatedAt != createdAt {
		return errors.New("error at create crontab")
	}

	return nil
}
