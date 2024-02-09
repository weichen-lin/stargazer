package main

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/consumer"
	database "github.com/weichen-lin/kafka-service/db"
)

func main() {

	driver, err := neo4j.NewDriverWithContext(
		"",
		neo4j.BasicAuth("neo4j", "", ""),
	)
	if err != nil {
		fmt.Println("Error creating driver:", err)
		return
	}

	pool, err := database.NewPostgresDB()
	if err != nil {
		fmt.Println("Error creating postgres connection:", err)
		return
	}

	get_user_consumer, err := consumer.GetUserProfileConsumer()
	get_repo_consumer, err := consumer.GetGithubReposConsumer()
	get_info_consumer, err := consumer.GetGithubRepoInfoConsumer()

	go get_user_consumer(driver)
	go get_repo_consumer(driver, pool)
	go get_info_consumer(pool)

	select {}
}
