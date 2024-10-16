package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/stargazer/domain"
)

var ErrNotFoundCrontab = errors.New("crontab not found")

func (db *Database) GetCrontab(ctx context.Context) (*domain.Crontab, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(db.Timeout)*time.Second)
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		session.Close(context.Background())
		cancel()
	}()

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User {email: $email})-[h:HAS_CRONTAB]-(c:Crontab)
			RETURN {
				triggered_at: c.triggered_at,
				created_at: c.created_at,
				updated_at: c.updated_at,
				status: c.status,
				last_triggered_at: c.last_triggered_at
			} as crontab
			`,
			map[string]interface{}{
				"email": email,
			})

		if err != nil {
			fmt.Println("error at read crontab: ", err)
			return nil, err
		}
		record, err := result.Single(ctx)
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
		return nil, fmt.Errorf("error converting crontab record to map")
	}

	cron, err := domain.FromCrontabEntity(
		&domain.CrontabEntity{
			TriggeredAt:     getString(crontab["triggered_at"]),
			CreatedAt:       getString(crontab["created_at"]),
			UpdatedAt:       getString(crontab["updated_at"]),
			Status:          getString(crontab["status"]),
			LastTriggeredAt: getString(crontab["last_triggered_at"]),
		},
	)

	if err != nil {
		return nil, err
	}

	return cron, nil
}

func (db *Database) SaveCrontab(ctx context.Context, crontab *domain.Crontab) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(db.Timeout)*time.Second)
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		session.Close(context.Background())
		cancel()
	}()

	entity := crontab.ToCrontabEntity()

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User {email: $email})
			MERGE (u)-[h:HAS_CRONTAB]-(c:Crontab)
			Set c += {
				created_at: $created_at,
				triggered_at: $triggered_at,
				updated_at: $updated_at,
				status: $status,
				last_triggered_at: $last_triggered_at
			}
			RETURN c.created_at AS created_at;
            `, map[string]any{
			"email":             email,
			"created_at":        entity.CreatedAt,
			"updated_at":        entity.UpdatedAt,
			"triggered_at":      entity.TriggeredAt,
			"last_triggered_at": entity.LastTriggeredAt,
			"status":            entity.Status,
		})

		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		return record, err
	})

	if err != nil {
		return err
	}

	record, ok := result.(*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting crontab records to *neo4j.Record")
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

type CrontabInfo struct {
	Email       string `json:"email"`
	TriggeredAt string `json:"triggered_at"`
}

func (db *Database) GetAllCrontab() []*CrontabInfo {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.Timeout)*time.Second)

	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		session.Close(context.Background())
		cancel()
	}()

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User)-[:HAS_CRONTAB]->(c:Crontab)
			RETURN u.email as email, c.triggered_at as triggered_at
			ORDER by c.created_at DESC
			`,
			map[string]interface{}{})

		if err != nil {
			return nil, err
		}

		record, err := result.Collect(ctx)
		return record, err
	})

	if err != nil {
		return nil
	}

	records, ok := result.([]*neo4j.Record)
	if !ok {
		return nil
	}

	crontabs := make([]*CrontabInfo, len(records))

	for i, record := range records {
		crontabs[i] = &CrontabInfo{
			Email:       getString(record.Values[0]),
			TriggeredAt: getString(record.Values[1]),
		}
	}

	return crontabs
}
