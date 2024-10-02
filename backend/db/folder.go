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
				id: $id,
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
				"id":         entity.Id,
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

func (db *Database) GetFolderById(ctx context.Context, id string) (*domain.Folder, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[h:HAS_FOLDER]-(f:Folder {id: $id})
			RETURN {
				id: f.id,
				name: f.name,
				is_public: f.is_public,
				created_at: f.created_at,
				updated_at: f.updated_at
			} as folder
			`,
			map[string]interface{}{
				"email": email,
				"id":    id,
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
			Id:        getString(data["id"]),
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

func (db *Database) GetFolderByName(ctx context.Context, name string) (*domain.Folder, error) {
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
				id: f.id,
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
			Id:        getString(data["id"]),
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

func (db *Database) DeleteFolder(ctx context.Context, id string) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	result, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[h:HAS_FOLDER]-(f:Folder {id: $id})
			OPTIONAL MATCH (r:Repository)-[i:IS_LOCATE]->(f)
			DELETE h, i, f
			RETURN elementId(f) as id
			`,
			map[string]interface{}{
				"email": email,
				"id":    id,
			})

		if err != nil {
			return nil, err
		}

		record, err := result.Collect(context.Background())
		return record, err
	})

	if err != nil {
		return err
	}

	folderRecord, ok := result.([]*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting folder records to *neo4j.Record")
	}

	for _, record := range folderRecord {
		record := record.AsMap()

		_, ok = record["id"].(string)
		if !ok {
			return fmt.Errorf("error convert id from record: %v", record)
		}
	}

	return nil
}

type PagingParams struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

type FolderSearchResult struct {
	Total int64                  `json:"total"`
	Data  []*domain.FolderEntity `json:"data"`
}

func (db *Database) GetFolders(ctx context.Context, params *PagingParams) (*FolderSearchResult, error) {
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
				id: f.id,
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
		return nil, fmt.Errorf("error at converting folder records to *neo4j.Record")
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
			Id:        getString(folderMap["id"]),
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

func (db *Database) AddRepoToFolder(ctx context.Context, folder *domain.Folder, repoIds []int64) error {
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
			MERGE (f:Folder {id: $id})
			ON CREATE SET f += {
				id: $id,
				name: $name,
				created_at: $created_at,
				updated_at: $updated_at
			}
			ON MATCH SET f.updated_at = $updated_at
			MERGE (u)-[:HAS_FOLDER]->(f)
			WITH u, f
			MATCH (r:Repository) 
			WHERE r.repo_id IN $repos
			MERGE (r)-[i:IS_LOCATE]->(f)
			ON MATCH SET i.created_at = datetime()
			ON CREATE SET i.created_at = datetime()
			RETURN i.created_at AS created_at
			`,
			map[string]interface{}{
				"email":      email,
				"repos":      repoIds,
				"id":         entity.Id,
				"name":       entity.Name,
				"created_at": entity.CreatedAt,
				"updated_at": entity.UpdatedAt,
			})

		if err != nil {
			return nil, err
		}
		record, err := result.Collect(context.Background())
		return record, err
	})

	if err != nil {
		return err
	}

	tagRecord, ok := result.([]*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting tag records to *neo4j.Record")
	}

	for _, record := range tagRecord {
		r := record.AsMap()
		_, ok = r["created_at"].(time.Time)
		if !ok {
			return fmt.Errorf("error convert name from record: %v", record)
		}
	}

	return nil
}

func (db *Database) DeleteRepoFromFolder(ctx context.Context, folder *domain.Folder, repoIds []int64) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	entity := folder.ToFolderEntity()

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	result, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[:HAS_FOLDER]->(f:Folder {id: $id})
			MATCH (r:Repository)-[i:IS_LOCATE]->(f)
			WHERE r.repo_id IN $repos
			SET f.updated_at = $updated_at
			DELETE i
			RETURN f.updated_at as updated_at
			`,
			map[string]interface{}{
				"email":      email,
				"id":         entity.Id,
				"repos":      repoIds,
				"updated_at": entity.UpdatedAt,
			})

		if err != nil {
			return nil, err
		}

		record, err := result.Collect(context.Background())
		if err != nil {
			return nil, err
		}

		return record, nil
	})

	if err != nil {
		return err
	}

	tagRecord, ok := result.([]*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting folder records to *neo4j.Record")
	}

	for _, record := range tagRecord {
		r := record.AsMap()
		_, ok = r["updated_at"].(string)
		if !ok {
			return fmt.Errorf("error convert name from record: %v", record)
		}
	}

	return nil
}

func (db *Database) GetFolderContainRepos(ctx context.Context, folder *domain.Folder, page int64, limit int64) (*SearchResult, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[h:HAS_FOLDER]-(f:Folder {id: $id})
			MATCH (r:Repository)-[i:IS_LOCATE]->(f)
			WITH f, COUNT(r) AS total
			MATCH (r:Repository)-[i:IS_LOCATE]->(f)
			WITH r, f, i.created_at AS created_at, total
			ORDER BY created_at DESC
			SKIP $limit * ($page - 1)
			LIMIT $limit
			RETURN total, collect({
				repo_id: r.repo_id,
				repo_name: r.repo_name,
				owner_name: r.owner_name,
				avatar_url: r.avatar_url,
				html_url: r.html_url,
				homepage: r.homepage,
				description: r.description,
				created_at: r.created_at,
				updated_at: r.updated_at,
				stargazers_count: r.stargazers_count,
				language: r.language,
				watchers_count: r.watchers_count,
				open_issues_count: r.open_issues_count,
				default_branch: r.default_branch,
				archived: r.archived,
				topics: r.topics,
				external_created_at: r.created_at,
				last_synced_at: r.last_synced_at,
				last_modified_at: r.last_modified_at
			}) as data
			`,
			map[string]interface{}{
				"email": email,
				"id":    folder.Id().String(),
				"limit": limit,
				"page":  page,
			})

		if err != nil {
			fmt.Println("error at read repo: ", err)
			return nil, err
		}
		record, err := result.Single(context.Background())
		return record, err
	})

	if err != nil {
		return &SearchResult{
			Data:  []*domain.RepositoryEntity{},
			Total: 0,
		}, nil
	}

	record, ok := result.(*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting users records to []*neo4j.Record")
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

	repos := make([]*domain.RepositoryEntity, len(data))

	for i, r := range data {
		repoMap := r.(map[string]interface{})

		entity := &domain.RepositoryEntity{
			RepoID:            getInt64(repoMap["repo_id"]),
			RepoName:          getString(repoMap["repo_name"]),
			OwnerName:         getString(repoMap["owner_name"]),
			AvatarURL:         getString(repoMap["avatar_url"]),
			HtmlURL:           getString(repoMap["html_url"]),
			Homepage:          getString(repoMap["homepage"]),
			Description:       getString(repoMap["description"]),
			CreatedAt:         getString(repoMap["created_at"]),
			UpdatedAt:         getString(repoMap["updated_at"]),
			StargazersCount:   getInt64(repoMap["stargazers_count"]),
			WatchersCount:     getInt64(repoMap["watchers_count"]),
			OpenIssuesCount:   getInt64(repoMap["open_issues_count"]),
			Language:          getString(repoMap["language"]),
			DefaultBranch:     getString(repoMap["default_branch"]),
			Archived:          getBool(repoMap["archived"]),
			Topics:            getStringArray(repoMap["topics"]),
			ExternalCreatedAt: getTimeString(repoMap["external_created_at"]),
			LastSyncedAt:      getTimeString(repoMap["last_synced_at"]),
			LastModifiedAt:    getTimeString(repoMap["last_modified_at"]),
		}

		repos[i] = entity
	}

	return &SearchResult{
		Data:  repos,
		Total: total,
	}, nil
}
