package db

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/domain"
)

func (db *Database) GetUser(ctx context.Context, email string) (*domain.User, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MATCH (u:User {email: $email})-[r]-(a:Account)
			RETURN {
				name: u.name,
				email: u.email,
				image: u.image,
				access_token: a.access_token,
				provider: a.provider,
				providerAccountId: a.providerAccountId,
				scope: a.scope,
				type: a.type,
				token_type: a.token_type 
			} as user
			`,
			map[string]interface{}{
				"email": email,
			})

		if err != nil {
			fmt.Println("error at read user: ", err)
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
		return nil, fmt.Errorf("error at converting crontab records to *neo4j.Record")
	}

	userMap := record.AsMap()

	userRecord, ok := userMap["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error converting record to map")
	}

	user := domain.FromUserEntity(
		&domain.UserEntity{
			Name:              getString(userRecord["name"]),
			Email:             getString(userRecord["email"]),
			Image:             getString(userRecord["image"]),
			AccessToken:       getString(userRecord["access_token"]),
			Provider:          getString(userRecord["provider"]),
			ProviderAccountId: getString(userRecord["providerAccountId"]),
			Scope:             getString(userRecord["scope"]),
			AuthType:          getString(userRecord["type"]),
			TokenType:         getString(userRecord["token_type"]),
		},
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
