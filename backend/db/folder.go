package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/stargazer/domain"
)

var ErrorNotFoundFolder = errors.New("folder not found")

func (db *Database) SaveFolder(ctx context.Context, folder *domain.Folder) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	entity := folder.ToFolderEntity()

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	result, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})
			MERGE (u)-[h:HAS_FOLDER]-(f:Folder {name: $name})
			ON CREATE SET f += {
				name: $name,
				is_public: $is_public,
				created_at: $created_at,
				updated_at: $updated_at
			}
			ON MATCH SET f += {
				name: $name,
				is_public: $is_public,
				updated_at: $updated_at
			}
			WITH f
			RETURN elementId(f) as id
			`,
			map[string]interface{}{
				"email":      email,
				"name":       entity.Name,
				"is_public":  entity.IsPublic,
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

	_, ok = record["id"].(string)
	if !ok {
		return fmt.Errorf("error convert name from record: %v", record)
	}

	return nil
}

func (db *Database) GetFolder(ctx context.Context, name string) (*domain.Folder, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[h:HAS_FOLDER]-(f:Folder {name: $name})
			RETURN {
				name: f.name,
				is_public: f.is_public,
				created_at: f.created_at,
				updated_at: f.updated_at
			} as folder
			`,
			map[string]interface{}{
				"email": email,
				"name":  name,
			})

		if err != nil {
			return nil, err
		}
		record, err := result.Single(context.Background())
		return record, err
	})

	if err != nil {
		return nil, err
	}

	folderRecord, ok := result.(*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting tag records to *neo4j.Record")
	}

	record := folderRecord.AsMap()

	data, ok := record["folder"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error convert name from record: %v", record)
	}

	folder, err := domain.FromFolderEntity(
		&domain.FolderEntity{
			Name:      getString(data["name"]),
			IsPublic:  getBool(data["is_public"]),
			CreatedAt: getString(data["created_at"]),
			UpdatedAt: getString(data["updated_at"]),
		},
	)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

type GetFolderParams struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

type FolderSearchResult struct {
	Total int64                  `json:"total"`
	Data  []*domain.FolderEntity `json:"data"`
}

func (db *Database) GetFolders(ctx context.Context, params *GetFolderParams) (*FolderSearchResult, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[h:HAS_FOLDER]-(f:Folder)
			WITH u, COUNT(f) as total
			MATCH (u)-[h:HAS_FOLDER]-(f)
			WITH total, h, f
			ORDER BY f.created_at DESC
			SKIP $limit * ($page - 1)
			LIMIT $limit
			RETURN total, collect({
				name: f.name,
				is_public: f.is_public,
				created_at: f.created_at,
				updated_at: f.updated_at
			}) as data
			`,
			map[string]interface{}{
				"email": email,
				"page":  params.Page,
				"limit": params.Limit,
			})

		if err != nil {
			fmt.Println("error at read folders: ", err)
			return nil, err
		}
		record, err := result.Single(context.Background())
		return record, err
	})

	if err != nil {
		return nil, err
	}

	record, ok := result.(*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting folder records to []*neo4j.Record")
	}

	recordMap := record.AsMap()

	total, ok := recordMap["total"].(int64)
	if !ok {
		return nil, fmt.Errorf("error convert id from record: %v", record)
	}

	data, ok := recordMap["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("error convert id from record: %v", record)
	}

	folders := make([]*domain.FolderEntity, len(data))

	for i, r := range data {
		folderMap := r.(map[string]interface{})

		entity := &domain.FolderEntity{
			Name:      getString(folderMap["name"]),
			IsPublic:  getBool(folderMap["is_public"]),
			CreatedAt: getString(folderMap["created_at"]),
			UpdatedAt: getString(folderMap["updated_at"]),
		}

		folders[i] = entity
	}

	return &FolderSearchResult{
		Data:  folders,
		Total: total,
	}, nil
}

func (db *Database) AddRepoToFolder(ctx context.Context, folder *domain.Folder, repo_id int64) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	entity := folder.ToFolderEntity()

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	result, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})
			MERGE (f:Folder {name: $name})
			ON CREATE SET f += {
				name: $name,
				created_at: $created_at,
				updated_at: $updated_at
			}
			ON MATCH SET f.updated_at = $updated_at
			MERGE (u)-[:HAS_FOLDER]->(f)
			WITH u, f
			MATCH (r:Repository {repo_id: $repo_id})
			MERGE (r)-[i:IS_LOCATE]->(f)
			ON MATCH SET i.created_at = datetime()
			ON CREATE SET i.created_at = datetime()
			RETURN i.created_at AS created_at
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
	_, ok = record["created_at"].(time.Time)
	if !ok {
		return fmt.Errorf("error convert name from record: %v", record)
	}

	return nil
}

func (db *Database) DeleteRepoFromFolder(ctx context.Context, folder *domain.Folder, repoID int64) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	entity := folder.ToFolderEntity()

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[:HAS_FOLDER]->(f:Folder {name: $name})
			MATCH (r:Repository {repo_id: $repo_id})-[i:IS_LOCATE]->(f)
			SET f.updated_at = $updated_at
			DELETE i
			RETURN f.updated_at
			`,
			map[string]interface{}{
				"email":      email,
				"name":       entity.Name,
				"repo_id":    repoID,
				"updated_at": entity.UpdatedAt,
			})

		if err != nil {
			return nil, err
		}

		record, err := result.Single(context.Background())
		if err != nil {
			return nil, err
		}

		return record, nil
	})

	if err != nil {
		return err
	}

	return nil
}
