package neo4j_kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/workflow"
)

func CreateUser(driver neo4j.DriverWithContext, user *workflow.User) (int64, error) {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	// Make constraint first
	constraint := `CREATE CONSTRAINT IF NOT EXISTS FOR (u:User) REQUIRE u.user_id IS UNIQUE`
	_, err := session.Run(context.Background(), constraint, nil)
	if err != nil {
		return 0, errors.New("error at create user constraint: " + err.Error())
	}

	user_id, err := session.ExecuteWrite(context.Background(), func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(context.Background(), `
			MERGE (u:User {user_id: $user_id})
			ON CREATE SET u.id = $id,
						u.name = $name,
						u.token = $token,
						u.is_sync = false,
						u.createdAt = datetime(),
						u.updatedAt = datetime()
			RETURN u.user_id AS user_id
			UNION
			MATCH (u:User {user_id: $user_id})
			RETURN u.user_id AS user_id;
            `,
			map[string]interface{}{
				"id":      uuid.New().String(),
				"user_id": user.ID,
				"name":    user.Login,
				"token":   "",
			})

		if err != nil {
			return 0, handleNeo4jError(err)
		}

		if result.Err() != nil {
			return 0, result.Err()
		}

		if result.Next(context.Background()) {
			record := result.Record()
			user_id, ok := record.Get("user_id")
			if !ok {
				return 0, fmt.Errorf("error at getting user_id from record: %v", record)
			}

			return user_id, nil
		}

		return 0, result.Err()
	})

	if user_id, ok := user_id.(int64); ok {
		return user_id, nil
	} else {
		return 0, fmt.Errorf("error at converting user_id to int64: %v", user_id)
	}
}

func CreateRepository(driver neo4j.DriverWithContext, repo *workflow.Repository, user_id int64) error {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	// Make constraint first
	constraint := `CREATE CONSTRAINT IF NOT EXISTS FOR (r:Repository) REQUIRE r.repo_id IS UNIQUE`
	_, err := session.Run(context.Background(), constraint, nil)
	if err != nil {
		return errors.New("error at create repo constraint: " + err.Error())
	}

	id, err := session.ExecuteWrite(context.Background(), func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(context.Background(), `
			MATCH (u:User {user_id: $user_id})
			MERGE (r:Repository {
			repo_id: $repo_id
			})
			ON CREATE SET
			r.id = $id,
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
			RETURN r.id AS id, r.repo_id AS repo_id, r.full_name AS full_name, r.default_branch AS default_branch;		
			`,
			map[string]interface{}{
				"user_id":          user_id,
				"id":               uuid.New().String(),
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
			id, ok := record.Get("id")
			if !ok {
				return nil, fmt.Errorf("repo_id not found")
			}
			return id, nil
		}

		return nil, result.Err()
	})

	if id, ok := id.(string); ok {
		return nil
	} else {
		return fmt.Errorf("error at converting repo_id to string: %v", id)
	}
}
