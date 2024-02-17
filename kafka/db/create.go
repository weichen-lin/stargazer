package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"unicode"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
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

type GetGithubReposInfo struct {
	UserId   string  `json:"user_id"`
	Username string `json:"username"`
	Page     int    `json:"page"`
}

type AddRepoInfo struct {
	Name          string `json:"name"`
	DefaultBranch string `json:"default_branch"`
	RepoId        int    `json:"id"`
	RepoInfo      []byte `json:"repo_info"`
}

func containsChinese(text string) bool {
	for _, char := range text {
		if unicode.Is(unicode.Scripts["Han"], char) {
			return true
		}
	}
	return false
}

func ReadMeTranslation(info string) (string, error) {
	token := os.Getenv("OPENAI_KEY")
	if token == "" {
		return "", fmt.Errorf("OPENAI_KEY not set")
	}

	client := openai.NewClient(token)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "This is a description of the Description information in a certain repository on GitHub. Please help me translate to English " + info,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("Error creating translation:", err)
	}

	return resp.Choices[0].Message.Content, nil
}

func ReadMeEmbedding(info string) ([]float32, error) {
	var vector []float32

	token := os.Getenv("OPENAI_KEY")
	if token == "" {
		return vector, fmt.Errorf("OPENAI_KEY not set")
	}
	client := openai.NewClient(token)

	// Create an EmbeddingRequest for the user query
	req := openai.EmbeddingRequest{
		Input: []string{info},
		Model: openai.AdaEmbeddingV2,
	}

	// Create an embedding for the user query
	resp, err := client.CreateEmbeddings(context.Background(), req)
	if err != nil {
		return vector, fmt.Errorf("Error creating query embedding:", err)
	}

	respEmbedding := resp.Data[0].Embedding

	if len(respEmbedding) == 0 {
		return vector, fmt.Errorf("Error creating query embedding: empty response")
	}

	return respEmbedding, nil
}

func CreateUser(driver neo4j.DriverWithContext, user *User) (string, error) {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	// Make constraint first
	constraint := `CREATE CONSTRAINT IF NOT EXISTS FOR (u:User) REQUIRE u.user_id IS UNIQUE`
	_, err := session.Run(context.Background(), constraint, nil)
	if err != nil {
		return "", errors.New("error at create user constraint: " + err.Error())
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

		err = handleNeo4jError(err)
		if err != nil {
			return 0, err
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

	if user_id, ok := user_id.(string); ok {
		return user_id, nil
	} else {
		return "", fmt.Errorf("error at converting user_id to int64: %v", user_id)
	}
}

func CreateRepository(driver neo4j.DriverWithContext, repo *Repository, user_id string, pool *gorm.DB) error {
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	tx := pool.Begin()
	if tx.Error != nil {
		return fmt.Errorf("error at begin transaction: %v", tx.Error)
	}

	go func() {
		var count int64
		if err := tx.Model(&RepoEmbeddingInfo{}).Where("repo_id = ?", repo.ID).Count(&count).Error; err != nil {
			tx.Rollback()
			return
		}

		if count > 0 {
			tx.Rollback()
			return
		}

		if err := tx.Create(&RepoEmbeddingInfo{
			RepoID:          repo.ID,
			FullName:        repo.FullName,
			Description:     repo.Description,
			AvatarURL:       repo.Owner.AvatarURL,
			HTMLURL:         repo.HTMLURL,
			StargazersCount: repo.StargazersCount,
			Language:        repo.Language,
			DefaultBranch:   repo.DefaultBranch,
		}).Error; err != nil {
			tx.Rollback()
			log.Fatalf("failed to insert repo data: %v. %d", err, repo.ID)
		} else {
			tx.Commit()
		}
	}()

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
