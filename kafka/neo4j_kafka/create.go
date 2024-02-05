package neo4j_kafka

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/workflow"
)

func CreateUser(driver neo4j.DriverWithContext, user workflow.User) error {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	// Make constraint first
	constraint := `CREATE CONSTRAINT IF NOT EXISTS FOR (u:User) REQUIRE u.user_id IS UNIQUE`
	_, err := session.Run(context.Background(), constraint, nil)
	fmt.Println("error at create constraint: ", err)

	_, err = session.ExecuteWrite(context.Background(), func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(context.Background(), `
			CREATE (u:User {
			id: apoc.create.uuid(),
			user_id: $user_id,
			name: $name,
			token: $token,
			is_sync: false,
			createdAt: datetime(),
			updatedAt: datetime()
			});
			`,
			map[string]interface{}{
				"user_id": user.ID,
				"name":    user.Login,
				"token":   "",
			})

		if err != nil {
			return nil, err
		}

		return nil, result.Err()
	})

	err = handleNeo4jError(err)
	return err
}

func CreateRepository(driver neo4j.DriverWithContext, repo workflow.Repository) error {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	// Make constraint first
	constraint := `CREATE CONSTRAINT IF NOT EXISTS FOR (r:Repository) REQUIRE r.repo_id IS UNIQUE`
	_, err := session.Run(context.Background(), constraint, nil)
	fmt.Println("error at create constraint: ", err)

	_, err = session.ExecuteWrite(context.Background(), func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(context.Background(), `
			CREATE (r:Repository {
			id: apoc.create.uuid(),
			repo_id: $repo_id,
			full_name: $full_name,
			avatar_url: $avatar_url,
			html_url: $html_url,
			description: $description,
			stargazers_count: $stargazers_count,
			language: $language,
			default_branch: $default_branch,
			last_updated_at: $last_updated_at,
			created_at: datetime()
			});
			`,
			map[string]interface{}{
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

		fmt.Println("error at create repo: ", err)

		if err != nil {
			return nil, err
		}

		return nil, result.Err()
	})

	err = handleNeo4jError(err)
	return err
}
