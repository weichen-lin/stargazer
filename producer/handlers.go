package main

import (
	"encoding/json"
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
	HandlerFunc func(db db.Database, msg *sarama.ConsumerMessage, producer sarama.SyncProducer) error
}

func GetGithubRepos(db db.Database, msg *sarama.ConsumerMessage, producer sarama.SyncProducer) error {
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

		user, err := db.GetUser(info.Email)

		params := &util.SendMailParams{
			Email:      info.Email,
			Name:       user.Name,
			StarsCount: (info.Page-1)*30 + len(stars),
		}

		err = util.SendMail(params)
		if err != nil {
			fmt.Println("Error sending email:", err)
		}
	}

	for _, repo := range stars {
		db.CreateRepository(&repo, info.Email)
	}

	return nil
}
