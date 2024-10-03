package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/stargazer/domain"
)

var ErrorNotFoundCollection = errors.New("collection not found")

func (db *Database) SaveCollection(ctx context.Context, collection *domain.Collection) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	entity := collection.ToCollectionEntity()

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	result, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})
			MERGE (u)-[h:HAS_COLLECT]-(c:Collection {name: $name})
			ON CREATE SET c += {
				id: $id,
				name: $name,
				is_public: $is_public,
				created_at: $created_at,
				updated_at: $updated_at
			}
			ON MATCH SET c += {
				name: $name,
				is_public: $is_public,
				updated_at: $updated_at
			}
			WITH c
			RETURN elementId(c) as id
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

	collectionRecord, ok := result.(*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting collection records to *neo4j.Record")
	}

	record := collectionRecord.AsMap()

	_, ok = record["id"].(string)
	if !ok {
		return fmt.Errorf("error convert name from record: %v", record)
	}

	return nil
}

func (db *Database) GetCollectionById(ctx context.Context, id string) (*domain.Collection, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[h:HAS_COLLECT]-(c:Collection {id: $id})
			RETURN {
				id: c.id,
				name: c.name,
				is_public: c.is_public,
				created_at: c.created_at,
				updated_at: c.updated_at
			} as collection
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

	collectionRecord, ok := result.(*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting tag records to *neo4j.Record")
	}

	record := collectionRecord.AsMap()

	data, ok := record["collection"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error convert name from record: %v", record)
	}

	collection, err := domain.FromCollectionEntity(
		&domain.CollectionEntity{
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

	return collection, nil
}

func (db *Database) GetCollectionByName(ctx context.Context, name string) (*domain.Collection, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[h:HAS_COLLECT]-(c:Collection {name: $name})
			RETURN {
				id: c.id,
				name: c.name,
				is_public: c.is_public,
				created_at: c.created_at,
				updated_at: c.updated_at
			} as collection
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

	collectionRecord, ok := result.(*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting collection records to *neo4j.Record")
	}

	record := collectionRecord.AsMap()

	data, ok := record["collection"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error convert name from record: %v", record)
	}

	collection, err := domain.FromCollectionEntity(
		&domain.CollectionEntity{
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

	return collection, nil
}

func (db *Database) DeleteCollection(ctx context.Context, id string) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	result, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[h:HAS_COLLECT]-(c:Collection {id: $id})
			OPTIONAL MATCH (r:Repository)-[i:IS_LOCATE]->(c)
			DELETE h, i, c
			RETURN elementId(c) as id
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

	collectionRecord, ok := result.([]*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting collection records to *neo4j.Record")
	}

	for _, record := range collectionRecord {
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

type CollectionSearchResult struct {
	Total int64                  `json:"total"`
	Data  []*domain.CollectionEntity `json:"data"`
}

func (db *Database) GetCollections(ctx context.Context, params *PagingParams) (*CollectionSearchResult, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[h:HAS_COLLECT]-(c:Collection)
			WITH u, COUNT(c) as total
			MATCH (u)-[h:HAS_COLLECT]-(c)
			WITH total, h, c
			ORDER BY c.created_at DESC
			SKIP $limit * ($page - 1)
			LIMIT $limit
			RETURN total, collect({
				id: c.id,
				name: c.name,
				is_public: c.is_public,
				created_at: c.created_at,
				updated_at: c.updated_at
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
		return nil, fmt.Errorf("error at converting collection records to *neo4j.Record")
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

	collections := make([]*domain.CollectionEntity, len(data))

	for i, r := range data {
		folderMap := r.(map[string]interface{})

		entity := &domain.CollectionEntity{
			Id:        getString(folderMap["id"]),
			Name:      getString(folderMap["name"]),
			IsPublic:  getBool(folderMap["is_public"]),
			CreatedAt: getString(folderMap["created_at"]),
			UpdatedAt: getString(folderMap["updated_at"]),
		}

		collections[i] = entity
	}

	return &CollectionSearchResult{
		Data:  collections,
		Total: total,
	}, nil
}

func (db *Database) AddRepoToCollection(ctx context.Context, collection *domain.Collection, repoIds []int64) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	entity := collection.ToCollectionEntity()

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	result, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})
			MERGE (c:Collection {id: $id})
			ON CREATE SET c += {
				id: $id,
				name: $name,
				created_at: $created_at,
				updated_at: $updated_at
			}
			ON MATCH SET c.updated_at = $updated_at
			MERGE (u)-[:HAS_COLLECT]->(c)
			WITH u, c
			MATCH (r:Repository) 
			WHERE r.repo_id IN $repos
			MERGE (r)-[i:IS_LOCATE]->(c)
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

	collectionRecord, ok := result.([]*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting collection records to *neo4j.Record")
	}

	for _, record := range collectionRecord {
		r := record.AsMap()
		_, ok = r["created_at"].(time.Time)
		if !ok {
			return fmt.Errorf("error convert name from record: %v", record)
		}
	}

	return nil
}

func (db *Database) DeleteRepoFromCollection(ctx context.Context, collection *domain.Collection, repoIds []int64) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	entity := collection.ToCollectionEntity()

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	result, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[:HAS_COLLECT]->(c:Collection {id: $id})
			MATCH (r:Repository)-[i:IS_LOCATE]->(c)
			WHERE r.repo_id IN $repos
			SET c.updated_at = $updated_at
			DELETE i
			RETURN c.updated_at as updated_at
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

	collectionRecord, ok := result.([]*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting collection records to *neo4j.Record")
	}

	for _, record := range collectionRecord {
		r := record.AsMap()
		_, ok = r["updated_at"].(string)
		if !ok {
			return fmt.Errorf("error convert name from record: %v", record)
		}
	}

	return nil
}

func (db *Database) GetCollectionContainRepos(ctx context.Context, collection *domain.Collection, page int64, limit int64) (*SearchResult, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[h:HAS_COLLECT]-(c:Collection {id: $id})
			MATCH (r:Repository)-[i:IS_LOCATE]->(c)
			WITH c, COUNT(r) AS total
			MATCH (r:Repository)-[i:IS_LOCATE]->(c)
			WITH r, c, i.created_at AS created_at, total
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
				"id":    collection.Id().String(),
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
