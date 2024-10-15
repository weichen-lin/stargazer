package db

import (
	"context"
	"errors"
	"net/mail"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/stargazer/domain"
)

var ErrNotFoundEmailAtContext = errors.New("not found email at context")

type EmailKey string

func (c EmailKey) String() string {
	return string(c)
}

func WithEmail(ctx context.Context, email string) (context.Context, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("invalid email format")
	}

	if len(email) < 5 || len(email) > 254 {
		return nil, errors.New("email length should be between 5 and 254 characters")
	}

	return context.WithValue(ctx, "email", email), nil
}

func EmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value("email").(string)
	return email, ok
}

func (db *Database) removeAllRecord() {
	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	_, _ = session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(context.Background(), `
			MATCH (n) DETACH DELETE n
			`,
			map[string]interface{}{})

		if err != nil {
			return nil, err
		}
		return nil, nil
	})
}

func getInt64(v interface{}) int64 {
	if i, ok := v.(int64); ok {
		return i
	}
	return 0
}

func getString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func getStringArray(v interface{}) []string {
	if s, ok := v.([]interface{}); !ok {
		return []string{}
	} else {
		var topics []string
		for _, topic := range s {
			topicString, ok := topic.(string)
			if !ok {
				continue
			}

			if topicString != "" {
				topics = append(topics, topicString)
			}
		}
		return topics
	}
}

func getBool(v interface{}) bool {
	if i, ok := v.(bool); ok {
		return i
	}
	return false
}

func getTimeString(v interface{}) string {
	if i, ok := v.(string); !ok {
		return time.Now().Format(time.RFC3339)
	} else {
		t, err := domain.ParseTime(i)

		if err != nil {
			return time.Now().Format(time.RFC3339)
		}

		return t.Format(time.RFC3339)
	}

}
