package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/segmentio/kafka-go"
	"github.com/weichen-lin/kafka-service/db"
	"github.com/weichen-lin/kafka-service/util"
)

var tokenCache = cache.New(20*time.Minute, 10*time.Minute)

type GetGithubReposInfo struct {
	Email string `json:"email"`
	Page  int    `json:"page"`
}

type RegisterConsumer struct {
	Topic       string
	HandlerFunc func(db *db.Database, msg kafka.Message, producer *kafka.Writer) error
}

func GetGithubRepos(db *db.Database, msg kafka.Message, producer *kafka.Writer) error {
	fmt.Printf("Received from Topic %s - Offset %d - Value : %s\n", msg.Topic, msg.Offset, msg.Value)

	var info GetGithubReposInfo

	err := json.Unmarshal(msg.Value, &info)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %s", err.Error())
	}

	var token string

	if tokenFormCache, found := tokenCache.Get(info.Email); !found {
		tokenValue, err := db.GetUserToken(info.Email)
		if err != nil {
			return fmt.Errorf("error getting user token %s", err.Error())
		}
		tokenCache.Set(info.Email, tokenValue, cache.DefaultExpiration)
		token = tokenValue
	} else {
		tokenValue, ok := tokenFormCache.(string)
		if !ok {
			return fmt.Errorf("error converting token to string: %s", token)
		}
		token = tokenValue
	}

	stars, err := util.GetUserStarredRepos(info.Page, token)

	if err != nil && errors.Is(err, util.ErrNoToken) {
		err = db.WriteResultAtCrontab(info.Email, "invalid github token")

		if err != nil {
			return fmt.Errorf("error writing result at crontab: %s", err.Error())
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("error getting user stars: %s", err.Error())
	}

	if len(stars) == 30 {
		info.Page++

		jsonString, err := json.Marshal(info)
		if err != nil {
			return fmt.Errorf("error marshalling JSON: %s", err.Error())
		}

		err = producer.WriteMessages(context.Background(), kafka.Message{
			Value: []byte(jsonString),
		})

		if err != nil {
			return fmt.Errorf("error sending message: %s", err.Error())
		}
	} else {
		starsCount := (info.Page-1)*30 + len(stars)

		err := util.SendMail(&util.SendMailParams{
			Email:      info.Email,
			Name:       "Stargazer user",
			StarsCount: starsCount,
		})

		err = db.WriteResultAtCrontab(info.Email, fmt.Sprintf("Successfully get %d starred repos", starsCount))
		if err != nil {
			return fmt.Errorf("error writing result at crontab: %s", err.Error())
		}

	}

	for _, repo := range stars {
		err := db.CreateRepository(&repo, info.Email)
		if err != nil {
			return fmt.Errorf("error creating repo: %s", err.Error())
		}
	}

	return nil
}
