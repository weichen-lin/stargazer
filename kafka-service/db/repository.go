package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/domain"
)

var (
	ErrRepositoryNotFound = errors.New("repository not found")
)

func (db *Database) GetRepository(ctx context.Context, repo_id int64) (*domain.Repository, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
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
			fmt.Println("error at read repo: ", err)
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
			StargazersCount: getInt(repo["stargazers_count"]),
			WatchersCount:   getInt(repo["watchers_count"]),
			OpenIssuesCount: getInt(repo["open_issues_count"]),
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

func (db *Database) CreateRepository(ctx context.Context, repo *domain.Repository) error {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return ErrNotFoundEmailAtContext
	}

	entity := repo.ToRepositoryEntity()
	
	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
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
			ON CREATE SET s += {
				is_delete : false,
				created_at : datetime()
			}
			ON MATCH SET s += {
				last_synced_at : datetime()
			}
			RETURN r.repo_id AS repo_id;
			`,
			map[string]interface{}{
				"email" : email,
				"repo_id": entity.RepoID,
				"repo_name": entity.RepoName,
				"owner_name": entity.OwnerName,
				"avatar_url": entity.AvatarURL,
				"html_url": entity.HtmlURL,
				"homepage": entity.Homepage,
				"description": entity.Description,
				"created_at": entity.CreatedAt,
				"updated_at": entity.UpdatedAt,
				"stargazers_count": entity.StargazersCount,
				"language": entity.Language,
				"watchers_count": entity.WatchersCount,
				"open_issues_count": entity.OpenIssuesCount,
				"default_branch": entity.DefaultBranch,
				"archived": entity.Archived,
			},
		)

		if err != nil {
			fmt.Println("error at create repo: ", err)
			return nil, err
		}
		records, err := result.Collect(context.Background())
		return records, err
	})

	if err != nil {
		return err
	}

	repos, ok := records.([]*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting users records to []*neo4j.Record")
	}

	if len(repos) == 0 {
		return fmt.Errorf("repo create failed")
	}

	record := repos[0].AsMap()
	_, ok = record["repo_id"].(int64)
	if !ok {
		return fmt.Errorf("error convert id from record: %v", record)
	}

	return nil
}
