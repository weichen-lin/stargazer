package neo4jOpeartion

import (
	"context"
	"errors"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type User struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
}

type Owner struct {
	Name 	string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Url 	string `json:"url"`
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

func CreateRepository(driver neo4j.DriverWithContext, repo *Repository, user_id string) error {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	// Make constraint first
	constraint := `CREATE CONSTRAINT IF NOT EXISTS FOR (r:Repository) REQUIRE r.repo_id IS UNIQUE`
	_, err := session.Run(context.Background(), constraint, nil)
	if err != nil {
		fmt.Println("error at create repo constraint: ", err)
		return errors.New("error at create repo constraint: " + err.Error())
	}

	id, err := session.ExecuteWrite(context.Background(), func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(context.Background(), `
			MATCH (u:User)-[re:HAS_ACCOUNT]-(a:Account { providerAccountId: $user_id })
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
			created_at: COALESCE(r.created_at, datetime()),
			open_issues_count: $open_issues_count,
			last_synced_at: datetime()
			}
			WITH u, r
			MERGE (u)-[s:STARS]->(r)
			MERGE (r)-[sb:STARRED_BY]->(u)
			SET s = {
			is_delete: COALESCE(s.is_delete, false)
			}
			RETURN r.repo_id AS repo_id;
			`,
			map[string]interface{}{
				"user_id":          user_id,
				"repo_id":          repo.ID,
				"full_name":        repo.FullName,
				"owner_name":       repo.Owner.Name,
				"owner_url":        repo.Owner.Url,
				"avatar_url":       repo.Owner.AvatarURL,
				"html_url":         repo.HTMLURL,
				"description":      repo.Description,
				"stargazers_count": repo.StargazersCount,
				"language":         repo.Language,
				"default_branch":   repo.DefaultBranch,
				"last_updated_at":  repo.UpdatedAt,
				"open_issues_count": repo.OpenIssuesCount,
			})

		if err != nil {
			fmt.Println("error at create repo: ", err)
			return nil, err
		}

		if result.Next(context.Background()) {
			record := result.Record()
			id, ok := record.Get("repo_id")
			if !ok {
				return nil, fmt.Errorf("repo_id not found")
			}
			return id, nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}

	if _, ok := id.(int64); ok {
		return nil
	} else {
		return fmt.Errorf("error at converting repo_id to string: %v", repo.ID)
	}
}


func AdjustCrontab(driver neo4j.DriverWithContext, hour int, user_name string) error {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	name, err := session.ExecuteWrite(context.Background(), func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(context.Background(), `
			MERGE (u:User {name: $name})-[h:HAS_CRONTAB]-(c:Crontab)
			SET c = {
				hour: $hour,
				updated_at: datetime(),
				created_at: COALESCE(c.created_at, datetime())
			}
			RETURN u.name AS name;
			`,
			map[string]interface{}{
				"name": user_name,
				"hour": hour,
			})

		if err != nil {
			fmt.Println("error at create crontab: ", err)
			return nil, err
		}

		if result.Next(context.Background()) {
			record := result.Record()
			_name, ok := record.Get("name")
			if !ok {
				return nil, fmt.Errorf("create failed")
			}
			return _name, nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}

	if _, ok := name.(string); ok {
		return nil
	} else {
		return fmt.Errorf("error at converting name to string: %v", user_name)
	}
}