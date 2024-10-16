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

	ctx, cancel := context.WithTimeout(ctx, time.Duration(db.Timeout)*time.Second)
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		session.Close(context.Background())
		cancel()
	}()

	entity := collection.ToCollectionEntity()

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User {email: $email})
			MERGE (u)-[h:HAS_COLLECT]-(c:Collection {id: $id})
			ON CREATE SET c += {
				id: $id,
				name: $name,
				description: $description,
				is_public: $is_public,
				created_at: $created_at,
				updated_at: $updated_at
			}
			ON MATCH SET c += {
				name: $name,
				description: $description,
				is_public: $is_public,
				updated_at: $updated_at
			}
			WITH c
			RETURN elementId(c) as id
			`,
			map[string]interface{}{
				"email":       email,
				"id":          entity.Id,
				"name":        entity.Name,
				"description": entity.Description,
				"is_public":   entity.IsPublic,
				"created_at":  entity.CreatedAt,
				"updated_at":  entity.UpdatedAt,
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

	collectionRecord, ok := result.(*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting collection record to *neo4j.Record")
	}

	record := collectionRecord.AsMap()

	_, ok = record["id"].(string)
	if !ok {
		return fmt.Errorf("error convert id from record: %v", record)
	}

	return nil
}

func (db *Database) GetCollectionByName(ctx context.Context, name string) (*domain.Collection, error) {
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
			MATCH (u:User {email: $email})-[h:HAS_COLLECT]-(c:Collection {name: $name})
			RETURN {
				id: c.id,
				name: c.name,
				description: c.description,
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
		record, err := result.Single(ctx)
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
			Id:          getString(data["id"]),
			Name:        getString(data["name"]),
			Description: getString(data["description"]),
			IsPublic:    getBool(data["is_public"]),
			CreatedAt:   getString(data["created_at"]),
			UpdatedAt:   getString(data["updated_at"]),
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

	ctx, cancel := context.WithTimeout(ctx, time.Duration(db.Timeout)*time.Second)
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		session.Close(context.Background())
		cancel()
	}()

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
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

		record, err := result.Collect(ctx)
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
	Total int64                      `json:"total"`
	Data  []*domain.CollectionEntity `json:"data"`
}

func (db *Database) GetCollections(ctx context.Context, params *PagingParams) (*CollectionSearchResult, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil,ErrNotFoundEmailAtContext
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(db.Timeout)*time.Second)
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		session.Close(context.Background())
		cancel()
	}()

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
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
				description: c.description,
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
			return nil, err
		}

		record, err := result.Single(ctx)
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
		return nil, fmt.Errorf("error convert total from record: %v", record)
	}

	data, ok := recordMap["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("error convert data from record: %v", record)
	}

	collections := make([]*domain.CollectionEntity, len(data))

	for i, r := range data {
		collectionMap := r.(map[string]interface{})

		entity := &domain.CollectionEntity{
			Id:          getString(collectionMap["id"]),
			Name:        getString(collectionMap["name"]),
			Description: getString(collectionMap["description"]),
			IsPublic:    getBool(collectionMap["is_public"]),
			CreatedAt:   getString(collectionMap["created_at"]),
			UpdatedAt:   getString(collectionMap["updated_at"]),
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

	ctx, cancel := context.WithTimeout(ctx, time.Duration(db.Timeout)*time.Second)
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		session.Close(context.Background())
		cancel()
	}()

	entity := collection.ToCollectionEntity()

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User {email: $email})
			MERGE (c:Collection {id: $id})
			ON CREATE SET c += {
				id: $id,
				name: $name,
				description: $description,
				created_at: $created_at,
				updated_at: $updated_at
			}
			ON MATCH SET c.updated_at = $updated_at
			MERGE (u)-[:HAS_COLLECT]->(c)
			WITH u, c
			MATCH (u)-[:STARS]-(r:Repository) 
			WHERE r.repo_id IN $repos
			MERGE (r)-[i:IS_LOCATE]->(c)
			ON MATCH SET i.updated_at = datetime()
			ON CREATE SET i.created_at = datetime()
			RETURN i.created_at AS created_at
			`,
			map[string]interface{}{
				"email":       email,
				"repos":       repoIds,
				"id":          entity.Id,
				"name":        entity.Name,
				"description": entity.Description,
				"created_at":  entity.CreatedAt,
				"updated_at":  entity.UpdatedAt,
			})

		if err != nil {
			return nil, err
		}
		record, err := result.Collect(ctx)
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
			return fmt.Errorf("error convert created_at from record: %v", record)
		}
	}

	return nil
}

func (db *Database) DeleteRepoFromCollection(ctx context.Context, collection *domain.Collection, repoIds []int64) error {
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

	entity := collection.ToCollectionEntity()

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
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

		record, err := result.Collect(ctx)
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

	ctx, cancel := context.WithTimeout(ctx, time.Duration(db.Timeout)*time.Second)
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		session.Close(context.Background())
		cancel()
	}()

	entity := collection.ToCollectionEntity()

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
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
				"id":    entity.Id,
				"limit": limit,
				"page":  page,
			})

		if err != nil {
			return nil, err
		}
		record, err := result.Single(ctx)
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
		return nil, fmt.Errorf("error at converting search records to *neo4j.Record")
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

func (db *Database) ShareCollection(ctx context.Context, collection *domain.Collection, shared_to string) error {
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

	entity := collection.ToCollectionEntity()

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User {email: $email})-[h:HAS_COLLECT]-(c:Collection {id: $id})
			MATCH (u2:User {email: $shared_to})
			MERGE (c)-[s:SHARED_WITH]->(u2)
			ON CREATE SET s.created_at = datetime()
			ON MATCH SET s.updated_at = datetime()
			RETURN elementId(s) as id
			`,
			map[string]interface{}{
				"email":     email,
				"id":        entity.Id,
				"shared_to": shared_to,
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

	collectionRecord, ok := result.(*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting collection records to *neo4j.Record")
	}

	record := collectionRecord.AsMap()

	_, ok = record["id"].(string)
	if !ok {
		return fmt.Errorf("error convert id from record: %v", record)
	}

	return nil
}

type SharedFrom struct {
	Email string `json:"email"`
	Image string `json:"image"`
	Name  string `json:"name"`
}

type SharedCollection struct {
	Owner      string                   `json:"owner"`
	Collection *domain.CollectionEntity `json:"collection"`
	SharedFrom *SharedFrom              `json:"shared_from"`
}

func (db *Database) GetCollectionById(ctx context.Context, id string) (*SharedCollection, error) {
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
			MATCH (c:Collection {id: $id})
			Optional MATCH (u2:User)-[:HAS_COLLECT]-(c)-[s:SHARED_WITH]->(u:User {email: $email})
			Optional MATCH (u3 {email: $email})-[:HAS_COLLECT]-(c)
			RETURN {
				id: c.id,
				name: c.name,
				description: c.description,
				is_public: c.is_public,
				created_at: c.created_at,
				updated_at: c.updated_at
			} as collection, {
				email: u2.email,
				image: u2.image,
				name: u2.name
			} as shared_from, {
				email: u3.email
			} as owner			 			
			`,
			map[string]interface{}{
				"email": email,
				"id":    id,
			})

		if err != nil {
			return nil, err
		}
		record, err := result.Single(ctx)
		return record, err
	})

	if err != nil {
		return nil, err
	}

	collectionRecord, ok := result.(*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting shared collection records to *neo4j.Record")
	}

	record := collectionRecord.AsMap()
	data, ok := record["collection"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error convert collection from record: %v", record)
	}

	shared := &SharedCollection{
		Collection: &domain.CollectionEntity{
			Id:          getString(data["id"]),
			Name:        getString(data["name"]),
			Description: getString(data["description"]),
			IsPublic:    getBool(data["is_public"]),
			CreatedAt:   getString(data["created_at"]),
			UpdatedAt:   getString(data["updated_at"]),
		},
		SharedFrom: nil,
		Owner:      getString(record["owner"].(map[string]interface{})["email"]),
	}

	user, ok := record["shared_from"].(map[string]interface{})
	if ok && getString(user["email"]) != "" {
		shared.SharedFrom = &SharedFrom{
			Email: getString(user["email"]),
			Image: getString(user["image"]),
			Name:  getString(user["name"]),
		}
	}

	return shared, nil
}

func (db *Database) GetSharedCollections(ctx context.Context) ([]*SharedCollection, error) {
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
			MATCH (u:User {email: $email})-[s:SHARED_WITH]-(c:Collection)
			MATCH (u2:User)-[h:HAS_COLLECT]-(c)
			RETURN {
				id: c.id,
				name: c.name,
				description: c.description,
				is_public: c.is_public,
				created_at: c.created_at,
				updated_at: c.updated_at
			} as collection, {
				email: u2.email,
				image: u2.image,
				name: u2.name 
			} as user
			`,
			map[string]interface{}{
				"email": email,
			})

		if err != nil {
			return nil, err
		}
		record, err := result.Collect(ctx)
		return record, err
	})

	if err != nil {
		return []*SharedCollection{}, nil
	}

	records, ok := result.([]*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting shared collection records to *neo4j.Record")
	}

	sharedCollection := make([]*SharedCollection, len(records))

	for i, record := range records {
		r := record.AsMap()

		collection, ok := r["collection"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("error convert name from record: %v", record)
		}

		user, ok := r["user"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("error convert user from record: %v", record)
		}

		sharedCollection[i] = &SharedCollection{
			Collection: &domain.CollectionEntity{
				Id:          getString(collection["id"]),
				Name:        getString(collection["name"]),
				Description: getString(collection["description"]),
				IsPublic:    getBool(collection["is_public"]),
				CreatedAt:   getString(collection["created_at"]),
				UpdatedAt:   getString(collection["updated_at"]),
			},
			SharedFrom: &SharedFrom{
				Email: getString(user["email"]),
				Image: getString(user["image"]),
				Name:  getString(user["name"]),
			},
		}
	}

	return sharedCollection, nil
}
