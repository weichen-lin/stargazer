package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/stargazer/domain"
)

var (
	ErrRepositoryNotFound = errors.New("repository not found")
)

func (db *Database) GetRepository(ctx context.Context, repo_id int64) (*domain.Repository, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[s:STARS {is_delete: false}]-(r:Repository  {repo_id: $repo_id})
			RETURN {
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
				external_created_at: s.created_at,
				last_synced_at: s.last_synced_at,
				last_modified_at: s.last_modified_at
			} as repo
			`,
			map[string]interface{}{
				"email":   email,
				"repo_id": repo_id,
			})

		if err != nil {
			return nil, err
		}

		record, err := result.Single(context.Background())
		return record, err
	})

	if err != nil {
		return nil, ErrRepositoryNotFound
	}

	record, ok := result.(*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting users records to *neo4j.Record")
	}

	repoMap := record.AsMap()

	repo, ok := repoMap["repo"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error converting record to map")
	}

	repository, err := domain.FromRepositoryEntity(
		&domain.RepositoryEntity{
			RepoID:            getInt64(repo["repo_id"]),
			RepoName:          getString(repo["repo_name"]),
			OwnerName:         getString(repo["owner_name"]),
			AvatarURL:         getString(repo["avatar_url"]),
			HtmlURL:           getString(repo["html_url"]),
			Homepage:          getString(repo["homepage"]),
			Description:       getString(repo["description"]),
			CreatedAt:         getString(repo["created_at"]),
			UpdatedAt:         getString(repo["updated_at"]),
			StargazersCount:   getInt64(repo["stargazers_count"]),
			WatchersCount:     getInt64(repo["watchers_count"]),
			OpenIssuesCount:   getInt64(repo["open_issues_count"]),
			Language:          getString(repo["language"]),
			DefaultBranch:     getString(repo["default_branch"]),
			Archived:          getBool(repo["archived"]),
			Topics:            getStringArray(repo["topics"]),
			ExternalCreatedAt: getTimeString(repo["external_created_at"]),
			LastSyncedAt:      getTimeString(repo["last_synced_at"]),
			LastModifiedAt:    getTimeString(repo["last_modified_at"]),
		},
	)

	if err != nil {
		return nil, err
	}

	return repository, nil
}

func (db *Database) DeleteRepository(ctx context.Context, repo_id int64) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	result, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[s:STARS]-(r:Repository  {repo_id: $repo_id})
			SET s.is_delete = true
			RETURN r.repo_id as repo_id
			`,
			map[string]interface{}{
				"email":   email,
				"repo_id": repo_id,
			})

		if err != nil {
			return nil, err
		}

		record, err := result.Single(context.Background())
		return record, err
	})

	if err != nil {
		fmt.Println(err)
		return ErrRepositoryNotFound
	}

	record, ok := result.(*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting users records to *neo4j.Record")
	}

	repoMap := record.AsMap()
	_, ok = repoMap["repo_id"].(int64)

	if !ok {
		return errors.New("failed to convert repo_id after delete")
	}

	return nil
}

type LanguageDistribution struct {
	Language string `json:"language"`
	Count    int64  `json:"count"`
}

func (db *Database) GetRepoLanguageDistribution(ctx context.Context) ([]*LanguageDistribution, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User { email: $email })-[s:STARS { is_delete: false }]->(r:Repository)
			WITH r.language as language, COUNT(r) as count
			RETURN language, count
			ORDER BY count DESC
			`,
			map[string]interface{}{
				"email": email,
			})

		if err != nil {
			fmt.Println("error at read repo: ", err)
			return nil, err
		}
		record, err := result.Collect(context.Background())
		return record, err
	})

	if err != nil {
		return nil, err
	}

	records, ok := result.([]*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting users records to []*neo4j.Record")
	}

	languages := make([]*LanguageDistribution, 0, len(records))

	for _, record := range records {
		languages = append(languages, &LanguageDistribution{
			Language: getString(record.Values[0]),
			Count:    getInt64(record.Values[1]),
		})
	}

	return languages, nil
}

func (db *Database) CreateRepository(ctx context.Context, repo *domain.Repository) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	entity := repo.ToRepositoryEntity()

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	records, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})
			MERGE (r:Repository { repo_id: $repo_id })
			SET r += {
				repo_id: $repo_id,
				repo_name: $repo_name,
				owner_name: $owner_name,
				avatar_url: $avatar_url,
				html_url: $html_url,
				homepage: $homepage,
				description: $description,
				created_at: $created_at,
				updated_at: $updated_at,
				stargazers_count: $stargazers_count,
				language: $language,
				watchers_count: $watchers_count,
				open_issues_count: $open_issues_count,
				default_branch: $default_branch,
				archived: $archived,
				topics: $topics
			}
			WITH u, r
			MERGE (u)-[s:STARS]->(r)
			MERGE (r)-[sb:STARRED_BY]->(u)
			ON CREATE 
			SET s += {
				is_delete : false,
				created_at : datetime()
			}
			ON MATCH 
			SET s += {
				last_synced_at : datetime()
			}
			WITH r, s, 
			CASE
				WHEN s.created_at = datetime() THEN 'CREATED'
				ELSE 'MATCHED'
			END AS operation
			RETURN 
			CASE operation
				WHEN 'CREATED' THEN s.created_at
				WHEN 'MATCHED' THEN s.last_synced_at
			END AS result
			`,
			map[string]interface{}{
				"email":             email,
				"repo_id":           entity.RepoID,
				"repo_name":         entity.RepoName,
				"owner_name":        entity.OwnerName,
				"avatar_url":        entity.AvatarURL,
				"html_url":          entity.HtmlURL,
				"homepage":          entity.Homepage,
				"description":       entity.Description,
				"created_at":        entity.CreatedAt,
				"updated_at":        entity.UpdatedAt,
				"stargazers_count":  entity.StargazersCount,
				"language":          entity.Language,
				"watchers_count":    entity.WatchersCount,
				"open_issues_count": entity.OpenIssuesCount,
				"default_branch":    entity.DefaultBranch,
				"archived":          entity.Archived,
				"topics":            entity.Topics,
			},
		)

		if err != nil {
			fmt.Println("error at create repo: ", err)
			return nil, err
		}
		records, err := result.Single(context.Background())
		return records, err
	})

	if err != nil {
		return ErrRepositoryNotFound
	}

	repos, ok := records.(*neo4j.Record)
	if !ok {
		return ErrRepositoryNotFound
	}

	record := repos.AsMap()
	_, ok = record["result"].(time.Time)
	if !ok {
		return fmt.Errorf("error convert id from record: %v", record)
	}

	return nil
}

type SearchParams struct {
	Page      int64    `json:"page"`
	Limit     int64    `json:"limit"`
	Languages []string `json:"languages"`
}

type SearchResult struct {
	Data  []*domain.RepositoryEntity `json:"data"`
	Total int64                      `json:"total"`
}

func (db *Database) SearchRepositoryByLanguage(ctx context.Context, params *SearchParams) (*SearchResult, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[s:STARS {is_delete: false}]-(r:Repository)
			WHERE r.language IN $languages
			WITH u, COUNT(r) as total
			MATCH (u)-[s:STARS]-(r)
			WHERE r.language IN $languages
			WITH total, s, r
			ORDER BY s.created_at DESC
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
				external_created_at: s.created_at,
				last_synced_at: s.last_synced_at,
				last_modified_at: s.last_modified_at
			}) as data
			`,
			map[string]interface{}{
				"email":     email,
				"languages": params.Languages,
				"limit":     params.Limit,
				"page":      params.Page,
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

type TopicResult struct {
	RepoId int64    `json:"repo_id"`
	Topics []string `json:"topics"`
}

func (db *Database) GetAllRepositoryTopics(ctx context.Context) ([]*TopicResult, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	results, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[s:STARS {is_delete: false}]-(r:Repository)
			RETURN r.repo_id as repo_id, r.topics as topics
			`,
			map[string]interface{}{
				"email": email,
			})

		if err != nil {
			fmt.Println("error at read repo: ", err)
			return nil, err
		}
		record, err := result.Collect(context.Background())
		return record, err
	})

	if err != nil {
		return nil, err
	}

	records, ok := results.([]*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting topics records to []*neo4j.Record")
	}

	topics := make([]*TopicResult, 0, len(records))
	for _, record := range records {
		topics = append(topics, &TopicResult{
			RepoId: getInt64(record.Values[0]),
			Topics: getStringArray(record.Values[1]),
		})
	}

	return topics, nil
}

type SortKey string

type SortOrder string

const (
	// StarGazer 建立時間
	SortKeyExternalCreatedAt SortKey = "created_at"

	// 星星數
	SortKeyStargazersCount SortKey = "stargazers_count"

	// 關注人數
	SortKeyWatchersCount SortKey = "watchers_count"
)

var sortKeyMap = map[SortKey]string{
	SortKeyExternalCreatedAt: "s.created_at",
	SortKeyStargazersCount:   "r.stargazers_count",
	SortKeyWatchersCount:     "r.watchers_count",
}

const (
	SortOrderDESC SortOrder = "DESC"
	SortOrderASC  SortOrder = "ASC"
)

var sortOrderMap = map[SortOrder]string{
	SortOrderDESC: "DESC",
	SortOrderASC:  "ASC",
}

type SortParams struct {
	Key   string
	Order string
}

var ErrInvalidSortKey = errors.New("invalid sort key")
var ErrInvalidSortOrder = errors.New("invalid sort order")

func (db *Database) GetRepositoriesOrderBy(ctx context.Context, params *SortParams) ([]*domain.RepositoryEntity, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	sortKey, exists := sortKeyMap[SortKey(params.Key)]
	if !exists {
		return []*domain.RepositoryEntity{}, ErrInvalidSortKey
	}

	sortOrder, exists := sortOrderMap[SortOrder(params.Order)]
	if !exists {
		return []*domain.RepositoryEntity{}, ErrInvalidSortOrder
	}

	query := fmt.Sprintf(`
			MATCH (u:User {email: $email})-[s:STARS { is_delete: false }]-(r:Repository)
			RETURN {
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
				external_created_at: s.created_at,
				last_synced_at: s.last_synced_at,
				last_modified_at: s.last_modified_at
			} as repo
			ORDER BY %s %s
			LIMIT 5
			`, sortKey, sortOrder)

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	results, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), query,
			map[string]interface{}{
				"email": email,
			})

		if err != nil {
			fmt.Println("error at read repo: ", err)
			return nil, err
		}
		record, err := result.Collect(context.Background())
		return record, err
	})

	if err != nil {
		return []*domain.RepositoryEntity{}, nil
	}

	records, ok := results.([]*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting repos records to []*neo4j.Record")
	}

	repos := make([]*domain.RepositoryEntity, len(records))

	for i, r := range records {
		record := r.AsMap()

		repoMap := record["repo"].(map[string]interface{})

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

	return repos, nil
}

func (db *Database) FullTextSearch(ctx context.Context, query string) ([]*domain.RepositoryEntity, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			CALL db.index.fulltext.queryNodes("REPOSITORY_FULL_TEXT_SEARCH", $query) YIELD node, score
			MATCH (User {
				email: $email
			})-[s:STARS]-(node)
			RETURN {
				repo_id: node.repo_id,
				repo_name: node.repo_name,
				owner_name: node.owner_name,
				avatar_url: node.avatar_url,
				html_url: node.html_url,
				homepage: node.homepage,
				description: node.description,
				created_at: node.created_at,
				updated_at: node.updated_at,
				stargazers_count: node.stargazers_count,
				language: node.language,
				watchers_count: node.watchers_count,
				open_issues_count: node.open_issues_count,
				default_branch: node.default_branch,
				archived: node.archived,
				topics: node.topics,
				external_created_at: s.created_at,
				last_synced_at: s.last_synced_at,
				last_modified_at: s.last_modified_at
			} AS data
			LIMIT 5
			`,
			map[string]interface{}{
				"email": email,
				"query": query,
			})

		if err != nil {
			fmt.Println("error at read repo: ", err)
			return nil, err
		}
		record, err := result.Collect(context.Background())
		return record, err
	})

	if err != nil {
		return []*domain.RepositoryEntity{}, nil
	}

	records, ok := result.([]*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting users records to []*neo4j.Record")
	}

	repos := make([]*domain.RepositoryEntity, len(records))

	for i, record := range records {
		recordMap := record.Values[0].(map[string]interface{})

		entity := &domain.RepositoryEntity{
			RepoID:            getInt64(recordMap["repo_id"]),
			RepoName:          getString(recordMap["repo_name"]),
			OwnerName:         getString(recordMap["owner_name"]),
			AvatarURL:         getString(recordMap["avatar_url"]),
			HtmlURL:           getString(recordMap["html_url"]),
			Homepage:          getString(recordMap["homepage"]),
			Description:       getString(recordMap["description"]),
			CreatedAt:         getString(recordMap["created_at"]),
			UpdatedAt:         getString(recordMap["updated_at"]),
			StargazersCount:   getInt64(recordMap["stargazers_count"]),
			WatchersCount:     getInt64(recordMap["watchers_count"]),
			OpenIssuesCount:   getInt64(recordMap["open_issues_count"]),
			Language:          getString(recordMap["language"]),
			DefaultBranch:     getString(recordMap["default_branch"]),
			Archived:          getBool(recordMap["archived"]),
			Topics:            getStringArray(recordMap["topics"]),
			ExternalCreatedAt: getTimeString(recordMap["external_created_at"]),
			LastSyncedAt:      getTimeString(recordMap["last_synced_at"]),
			LastModifiedAt:    getTimeString(recordMap["last_modified_at"]),
		}

		repos[i] = entity
	}

	return repos, nil
}
