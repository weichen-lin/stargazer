package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/domain"
)

var ErrNotFoundUser = errors.New("user not found")

func (db *Database) CreateUser(user *domain.User) error {
	entity := user.ToUserEntity()

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	result, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(context.Background(), `
			MERGE (u:User {email: $email})
			ON CREATE SET u.name = $name,
						u.image = $image
			ON MATCH SET u.name = $name,
						u.image = $image
			MERGE (a:Account {providerAccountId: $provider_account_id})
			ON CREATE SET a.access_token = $access_token,
						a.provider = $provider,
						a.scope = $scope,
						a.type = $auth_type,
						a.token_type = $token_type
			ON MATCH SET a.access_token = $access_token,
						a.provider = $provider,
						a.scope = $scope,
						a.type = $auth_type,
						a.token_type = $token_type
			MERGE (u)-[:HAS_ACCOUNT]->(a)
			RETURN u.name AS name;
			`,
			map[string]interface{}{
				"name":                entity.Name,
				"email":               entity.Email,
				"image":               entity.Image,
				"access_token":        entity.AccessToken,
				"provider":            entity.Provider,
				"provider_account_id": entity.ProviderAccountId,
				"scope":               entity.Scope,
				"auth_type":           entity.AuthType,
				"token_type":          entity.TokenType,
			})

		if err != nil {
			fmt.Println("error at create user: ", err)
			return nil, err
		}
		record, err := result.Single(context.Background())
		return record, err
	})

	if err != nil {
		return err
	}

	userRecord, ok := result.(*neo4j.Record)
	if !ok {
		return fmt.Errorf("error at converting users records to []*neo4j.Record")
	}

	record := userRecord.AsMap()
	_, ok = record["name"].(string)
	if !ok {
		return fmt.Errorf("error convert name from record: %v", record)
	}

	return nil
}

func (db *Database) GetUser(ctx context.Context) (*domain.User, error) {
	email, ok := EmailFromContext(ctx)
	if !ok {
		return nil, ErrNotFoundEmailAtContext
	}

	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
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
		return nil, ErrNotFoundUser
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

	return user, nil
}

func (db *Database) DeleteUser(ctx context.Context) error {
	return nil
}
