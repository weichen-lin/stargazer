package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image"`
}

type Owner struct {
	Name      string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Url       string `json:"url"`
}

type Repository struct {
	ID              int64  `json:"id"`
	FullName        string `json:"full_name"`
	Owner           Owner  `json:"owner"`
	HTMLURL         string `json:"html_url"`
	Description     string `json:"description"`
	UpdatedAt       string `json:"updated_at"`
	StargazersCount int    `json:"stargazers_count"`
	Language        string `json:"language"`
	DefaultBranch   string `json:"default_branch"`
	OpenIssuesCount int    `json:"open_issues_count"`
}

func (db *Database) CreateRepository(repo *Repository, email string) error {
	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	// Make constraint first
	constraint := `CREATE CONSTRAINT IF NOT EXISTS FOR (r:Repository) REQUIRE r.repo_id IS UNIQUE`
	_, err := session.Run(context.Background(), constraint, nil)
	if err != nil {
		fmt.Println("error at create repo constraint: ", err)
		return errors.New("error at create repo constraint: " + err.Error())
	}

	records, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})
			MERGE (r:Repository { repo_id: $repo_id })
			SET r = {
				repo_id: $repo_id,
				full_name: $full_name,
				owner_name: $owner_name,
				owner_url: $owner_url,
				avatar_url: $avatar_url,
				html_url: $html_url,
				description: $description,
				stargazers_count: $stargazers_count,
				language: $language,
				default_branch: $default_branch,
				last_updated_at: $last_updated_at,
				open_issues_count: $open_issues_count
			}
			WITH u, r
			MERGE (u)-[s:STARS]->(r)
			MERGE (r)-[sb:STARRED_BY]->(u)
			SET s += {
				is_delete: COALESCE(s.is_delete, false),
				created_at: COALESCE(r.created_at, datetime()),
				last_synced_at: datetime()
			}
			RETURN r.repo_id AS repo_id;
			`,
			map[string]interface{}{
				"email":             email,
				"repo_id":           repo.ID,
				"full_name":         repo.FullName,
				"owner_name":        repo.Owner.Name,
				"owner_url":         repo.Owner.Url,
				"avatar_url":        repo.Owner.AvatarURL,
				"html_url":          repo.HTMLURL,
				"description":       repo.Description,
				"stargazers_count":  repo.StargazersCount,
				"language":          repo.Language,
				"default_branch":    repo.DefaultBranch,
				"last_updated_at":   repo.UpdatedAt,
				"open_issues_count": repo.OpenIssuesCount,
			})

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

func (db *Database) UpdateCrontab(hour int, email string) error {
	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	records, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})
			MERGE (u)-[h:HAS_CRONTAB]-(c:Crontab)
			SET c += {
				hour: $hour,
				updated_at: datetime(),
				created_at: COALESCE(c.created_at, datetime())
			}
			RETURN c.updated_at AS updated_at;
            `, map[string]any{
			"email": email,
			"hour":  hour,
		})

		if err != nil {
			return nil, err
		}

		records, _ := result.Collect(context.Background())
		if err != nil {
			return nil, err
		}

		return records, nil
	})

	if err != nil {
		return err
	}

	users, ok := records.([]*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting users records to []*neo4j.Record")
	}

	if len(users) == 0 {
		return fmt.Errorf("user not found")
	}

	record := users[0].AsMap()
	_, ok = record["updated_at"].(time.Time)

	if !ok {
		return fmt.Errorf("error at convert update time from record: %v", record)
	}

	return nil
}
