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
			MATCH (u:User {email: $email})-[s:STARS]-(r:Repository  {repo_id: $repo_id})
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
				archived: r.archived
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
		return nil, fmt.Errorf("error at converting users records to []*neo4j.Record")
	}

	repoMap := record.AsMap()

	repo, ok := repoMap["repo"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error converting record to map")
	}

	repository, err := domain.FromRepositoryEntity(
		&domain.RepositoryEntity{
			RepoID:          getInt64(repo["repo_id"]),
			RepoName:        getString(repo["repo_name"]),
			OwnerName:       getString(repo["owner_name"]),
			AvatarURL:       getString(repo["avatar_url"]),
			HtmlURL:         getString(repo["html_url"]),
			Homepage:        getString(repo["homepage"]),
			Description:     getString(repo["description"]),
			CreatedAt:       getString(repo["created_at"]),
			UpdatedAt:       getString(repo["updated_at"]),
			StargazersCount: getInt64(repo["stargazers_count"]),
			WatchersCount:   getInt64(repo["watchers_count"]),
			OpenIssuesCount: getInt64(repo["open_issues_count"]),
			Language:        getString(repo["language"]),
			DefaultBranch:   getString(repo["default_branch"]),
			Archived:        getBool(repo["archived"]),
		},
	)

	if err != nil {
		return nil, err
	}

	return repository, nil
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
				archived: $archived
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
				archived: r.archived
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
			RepoID:          getInt64(repoMap["repo_id"]),
			RepoName:        getString(repoMap["repo_name"]),
			OwnerName:       getString(repoMap["owner_name"]),
			AvatarURL:       getString(repoMap["avatar_url"]),
			HtmlURL:         getString(repoMap["html_url"]),
			Homepage:        getString(repoMap["homepage"]),
			Description:     getString(repoMap["description"]),
			CreatedAt:       getString(repoMap["created_at"]),
			UpdatedAt:       getString(repoMap["updated_at"]),
			StargazersCount: getInt64(repoMap["stargazers_count"]),
			WatchersCount:   getInt64(repoMap["watchers_count"]),
			OpenIssuesCount: getInt64(repoMap["open_issues_count"]),
			Language:        getString(repoMap["language"]),
			DefaultBranch:   getString(repoMap["default_branch"]),
			Archived:        getBool(repoMap["archived"]),
		}

		repos[i] = entity
	}

	return &SearchResult{
		Data:  repos,
		Total: total,
	}, nil
}