package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Config struct {
	OpenAiKey   string  `json:"openai_key"`
	GithubToken string  `json:"github_token"`
	Limit       int64   `json:"limit"`
	Cosine      float64 `json:"cosine"`
}

type Crontab struct {
	Email    string    `json:"email"`
	Note     string    `json:"note"`
	Hour     int64     `json:"hour"`
	Status   string    `json:"status"`
	UpdateAt time.Time `json:"update_at"`
}

func (db *database) GetUser(email string) (*User, error) {
	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	records, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, _ := tx.Run(context.Background(), `
            MATCH (u:User {email: $email}) 
            RETURN u.id as id, u.name as name, u.email as email, u.image as image
            `, map[string]any{
			"email": email,
		})
		records, _ := result.Collect(context.Background())
		return records, nil
	})

	if err != nil {
		return nil, err
	}

	users, ok := records.([]*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting users records to []*neo4j.Record")
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	record := users[0].AsMap()
	id, ok := record["id"].(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("error at getting id from record: %v", record)
	}

	name, ok := record["name"].(string)
	if !ok {
		return nil, fmt.Errorf("error at getting name from record: %v", record)
	}

	image, ok := record["image"].(string)
	if !ok {
		return nil, fmt.Errorf("error at getting image from record: %v", record)
	}

	return &User{
		ID:    id,
		Name:  name,
		Email: email,
		Image: image,
	}, nil
}

func (db *database) GetUserConfig(email string) (*Config, error) {
	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	records, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, _ := tx.Run(context.Background(), `
            MATCH (u:User {email: $email})-[:HAS_CONFIG]->(c:Config)
            RETURN c.openai_key as openai_key, c.github_token as github_token, c.limit as limit, c.cosine as cosine
            `, map[string]any{
			"email": email,
		})
		records, _ := result.Collect(context.Background())
		return records, nil
	})

	if err != nil {
		return nil, err
	}

	configs, ok := records.([]*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting users records to []*neo4j.Record")
	}

	if len(configs) == 0 {
		return nil, fmt.Errorf("config not setting")
	}

	record := configs[0].AsMap()
	openai_key, ok := record["openai_key"].(string)
	if !ok {
		return nil, fmt.Errorf("error at convert openai_key from record: %v", record)
	}

	github_token, ok := record["github_token"].(string)
	if !ok {
		return nil, fmt.Errorf("error at convert openai_key from record: %v", record)
	}

	limit, ok := record["limit"].(int64)
	if !ok {
		return nil, fmt.Errorf("error at convert limit from record: %v", record)
	}

	cosine, ok := record["cosine"].(float64)
	if !ok {
		return nil, fmt.Errorf("error at convert cosine from record: %v", record)
	}

	return &Config{
		OpenAiKey:   openai_key,
		GithubToken: github_token,
		Limit:       limit,
		Cosine:      cosine,
	}, nil
}

func (db *database) GetUserNotVectorize(userName string) ([]int64, error) {
	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	records, err := session.ExecuteRead(context.Background(), func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(context.Background(), `
				MATCH (u:User {name: $name})-[s:STARS]-(r:Repository)
				WHERE s.is_vectorized = FALSE or s.is_vectorized IS NULL
				RETURN r.repo_id as repo_id
            `,
			map[string]interface{}{
				"name": userName,
			})

		if err != nil {
			return nil, err
		}

		if result.Err() != nil {
			return nil, result.Err()
		}

		collects, err := result.Collect(context.Background())
		if err != nil {
			return nil, err
		}

		stars := make([]int64, len(collects))

		for i, record := range collects {
			repo_id, ok := record.Get("repo_id")
			if !ok {
				return nil, fmt.Errorf("error at getting repo_id from record: %v", record)
			}

			stars[i] = repo_id.(int64)
		}

		return stars, result.Err()
	})

	if err != nil {
		fmt.Println("Error make vectorize success from neo4j:", err)
		return nil, err
	}

	if _, ok := records.([]int64); !ok {
		return nil, fmt.Errorf("error at converting stars to []int")
	}

	return records.([]int64), err
}

func (db *database) GetAllUserCrontab() ([]Crontab, error) {
	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	records, err := session.ExecuteRead(context.Background(), func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, _ := transaction.Run(context.Background(), `
				MATCH (u:User)-[h:HAS_CRONTAB]-(c:Crontab)
				RETURN u.email as email, c.hour as hour            
				`,
			map[string]interface{}{})

		records, _ := result.Collect(context.Background())
		return records, nil
	})

	if err != nil {
		return nil, err
	}

	crontabs, ok := records.([]*neo4j.Record)
	if !ok {
		return nil, fmt.Errorf("error at converting users records to []*neo4j.Record")
	}

	if len(crontabs) == 0 {
		return []Crontab{}, nil
	}

	crontabList := make([]Crontab, len(crontabs))

	for i, record := range crontabs {
		record := record.AsMap()
		email, ok := record["email"].(string)
		if !ok {
			return nil, fmt.Errorf("error at getting email from record: %v", record)
		}

		hour, ok := record["hour"].(int64)
		if !ok {
			return nil, fmt.Errorf("error at getting hour from record: %v", record)
		}

		crontabList[i] = Crontab{
			Email: email,
			Hour:  hour,
		}
	}

	return crontabList, nil
}
