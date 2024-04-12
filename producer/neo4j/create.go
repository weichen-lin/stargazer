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
	AvatarURL string `json:"avatar_url"`
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
			MERGE (r:Repository {
			repo_id: $repo_id
			})
			ON CREATE SET
			r.repo_id = $repo_id,
			r.full_name = $full_name,
			r.avatar_url = $avatar_url,
			r.html_url = $html_url,
			r.description = $description,
			r.stargazers_count = $stargazers_count,
			r.language = $language,
			r.default_branch = $default_branch,
			r.last_updated_at = $last_updated_at,
			r.created_at = datetime()
			WITH u, r
			MERGE (u)-[s:STARS]->(r)
			MERGE (r)-[sb:STARRED_BY]->(u)
			RETURN r.repo_id AS repo_id;	
			`,
			map[string]interface{}{
				"user_id":          user_id,
				"repo_id":          repo.ID,
				"full_name":        repo.FullName,
				"avatar_url":       repo.Owner.AvatarURL,
				"html_url":         repo.HTMLURL,
				"description":      repo.Description,
				"stargazers_count": repo.StargazersCount,
				"language":         repo.Language,
				"default_branch":   repo.DefaultBranch,
				"last_updated_at":  repo.UpdatedAt,
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
