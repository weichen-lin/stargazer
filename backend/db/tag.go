package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/stargazer/domain"
)

type CreateTagPayload struct {
	Name   string
	RepoID int64
}

var ErrorNotFoundTag = errors.New("tag not found")

func (db *Database) SaveTag(ctx context.Context, tag *domain.Tag, repo_id int64) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	entity := tag.ToTagEntity()

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	result, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})
			MERGE (t:Tag {name: $name})
			ON CREATE SET t += {
				name: $name,
				created_at: $created_at,
				updated_at: $updated_at
			}
			ON MATCH SET t.updated_at = $updated_at
			MERGE (u)-[:HAS_TAG]->(t)
			WITH u, t
			MATCH (r:Repository {repo_id: $repo_id})
			MERGE (u)-[s:STARS]->(r)
			ON MATCH SET s.last_modified_at = datetime()
			ON CREATE SET s.last_modified_at = datetime()
			MERGE (r)-[:TAGGED_BY]->(t)
			RETURN t.created_at AS created_at
			`,
			map[string]interface{}{
				"email":      email,
				"repo_id":    repo_id,
				"name":       entity.Name,
				"created_at": entity.CreatedAt,
				"updated_at": entity.UpdatedAt,
			})

		if err != nil {
			return nil, err
		}
		record, err := result.Single(context.Background())
		return record, err
	})

	if err != nil {
		return err
	}

	tagRecord, ok := result.(*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting tag records to *neo4j.Record")
	}

	record := tagRecord.AsMap()
	_, ok = record["created_at"].(string)
	if !ok {
		return fmt.Errorf("error convert name from record: %v", record)
	}

	return nil
}

func (db *Database) RemoveTag(ctx context.Context, tag *domain.Tag, repoID int64) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
            MATCH (u:User {email: $email})-[h:HAS_TAG]->(t:Tag {name: $name})
            MATCH (t)<-[tb:TAGGED_BY]-(r:Repository {repo_id: $repo_id})
            DELETE h, t, tb
            `,
			map[string]interface{}{
				"email":   email,
				"name":    tag.Name(),
				"repo_id": repoID,
			})

		if err != nil {
			return nil, err
		}

		return result, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetTagByName(ctx context.Context, name string) (*domain.Tag, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[r:HAS_TAG]-(t:Tag {name: $name})
			RETURN {
				name: t.name,
				created_at: t.created_at,
				updated_at: t.updated_at
			} as tag
			`,
			map[string]interface{}{
				"email": email,
				"name":  name,
			})

		if err != nil {
			fmt.Println("error at get tag: ", err)
			return nil, err
		}
		record, err := result.Single(context.Background())
		return record, err
	})

	if err != nil {
		return nil, ErrorNotFoundTag
	}

	record, ok := result.(*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting tag record to *neo4j.Record")
	}

	tagMap := record.AsMap()
	tagData, ok := tagMap["tag"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error converting record to map")
	}

	tag, err := domain.FromTagEntity(
		&domain.TagEntity{
			Name:      getString(tagData["name"]),
			CreatedAt: getString(tagData["created_at"]),
			UpdatedAt: getString(tagData["updated_at"]),
		},
	)

	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (db *Database) GetTagsByRepo(ctx context.Context, repo_id int64) ([]*domain.Tag, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[:HAS_TAG]->(t:Tag)<-[:TAGGED_BY]-(r:Repository {repo_id: $repo_id})
  			RETURN t.created_at as created_at, t.updated_at as updated_at, t.name as name
			`,
			map[string]interface{}{
				"email":   email,
				"repo_id": repo_id,
			})

		if err != nil {
			return nil, err
		}
		record, err := result.Collect(context.Background())
		return record, err
	})

	if err != nil {
		return []*domain.Tag{}, err
	}

	records, ok := result.([]*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting users records to []*neo4j.Record")
	}

	tags := make([]*domain.Tag, 0, len(records))

	for _, record := range records {
		tag, err := domain.FromTagEntity(
			&domain.TagEntity{
				Name:      getString(record.Values[2]),
				CreatedAt: getString(record.Values[0]),
				UpdatedAt: getString(record.Values[1]),
			},
		)

		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
