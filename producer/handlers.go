package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/patrickmn/go-cache"
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
	HandlerFunc func(db *db.Database, msg *sarama.ConsumerMessage, producer sarama.SyncProducer) error
}

func GetGithubRepos(db *db.Database, msg *sarama.ConsumerMessage, producer sarama.SyncProducer) error {
	fmt.Printf("Received message: Topic - %s, Key - %s, Value - %s\n",
		msg.Topic, msg.Key, msg.Value)

	var info GetGithubReposInfo

	err := json.Unmarshal(msg.Value, &info)
	if err != nil {
		return err
	}

	var token string

	if tokenFormCache, found := tokenCache.Get(info.Email); !found {
		tokenValue, err := db.GetUserConfig(info.Email)
		if err != nil {
			return fmt.Errorf("Error getting user token %s:", err.Error())
		}
		token = tokenValue.GithubToken
		tokenCache.Set(info.Email, tokenValue.GithubToken, cache.DefaultExpiration)
	} else {
		tokenValue, ok := tokenFormCache.(string)
		if !ok {
			return fmt.Errorf("Error converting token to string: %s", token)
		}
		token = tokenValue
	}

	stars, err := util.GetUserStarredRepos(info.Page, token)
	if err != nil && errors.Is(err, util.ErrNoToken) {
		err = db.WriteResultAtCrontab(info.Email, "invalid github token")
		
		if err != nil {
			return fmt.Errorf("Error writing result at crontab: %s", err.Error())
		}
	}

	if err != nil {
		return fmt.Errorf("Error getting user stars: %s", err.Error())
	}

	if len(stars) == 30 {
		info.Page++

		jsonString, err := json.Marshal(info)
		if err != nil {
			return fmt.Errorf("Error marshalling JSON: %s", err.Error())
		}

		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: "get_user_stars",
			Value: sarama.StringEncoder(jsonString),
		})

		if err != nil {
			return fmt.Errorf("Error sending message: %s", err.Error())
		}
	} else {
		
		err = db.WriteResultAtCrontab(info.Email, fmt.Sprintf("Successfully get %d starred repos", (info.Page-1)*30+len(stars)))
		if err != nil {
			return fmt.Errorf("Error writing result at crontab: %s", err.Error())
		}

		user, err := db.GetUser(info.Email)
		if err != nil {
			return fmt.Errorf("Error getting user: %s", err.Error())
		}

		params := &util.SendMailParams{
			Email:      info.Email,
			Name:       user.Name,
			StarsCount: (info.Page-1)*30 + len(stars),
		}

		err = util.SendMail(params)
		if err != nil {
			return fmt.Errorf("Error sending mail: %s", err.Error())
		}
	}

	for _, repo := range stars {
		err := db.CreateRepository(&repo, info.Email)
		if err != nil {
			return fmt.Errorf("Error creating repo: %s", err.Error())
		}
	}

	return nil
}
