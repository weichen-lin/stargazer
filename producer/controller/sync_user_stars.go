package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/kafka-service/db"
	"github.com/weichen-lin/kafka-service/workflow"
)

type Client struct {
	Email      string
	StatusCode chan int
}

type SyncUserStars struct {
	Stars []int  `json:"stars" binding:"required" form:"stars"`
	Email string `json:"email" binding:"required" form:"email"`
}

var clients = make(map[string]*Client)

func HandleConnections(c *gin.Context) {

	email, ok := c.Value("email").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	driver, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db not found"})
		return
	}

	driver, ok = driver.(db.Database)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Neo4j driver"})
		return
	}

	stars, err := driver.(db.Database).GetUserNotVectorize(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user stars"})
		return
	}

	client := &Client{
		Email:      email,
		StatusCode: make(chan int),
	}

	clients[email] = client
	defer delete(clients, email)

	c.Stream(func(w io.Writer) bool {
		w.(http.ResponseWriter).Header().Set("Content-Type", "text/event-stream")
		w.(http.ResponseWriter).Header().Set("Cache-Control", "no-cache")
		w.(http.ResponseWriter).Header().Set("Connection", "keep-alive")
		w.(http.ResponseWriter).Header().Set("Access-Control-Allow-Origin", "*")

		c.SSEvent("message", map[string]interface{}{"total": len(stars)})
		c.Writer.Flush()

		go func(driver db.Database) {
			defer close(client.StatusCode)
			for i := 0; i < len(stars); i++ {

				id := stars[i]

				_, err := workflow.VectorizeStar(&workflow.SyncUserStar{
					Email:  email,
					RepoId: id,
				})

				if err != nil {
					c.SSEvent("message", map[string]interface{}{"error": err.Error()})
					c.Writer.Flush()
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				client.StatusCode <- i

				err = driver.(db.Database).ConfirmVectorize(&workflow.SyncUserStar{
					Email:  email,
					RepoId: id,
				})

				if err != nil {
					c.SSEvent("message", map[string]interface{}{"error": err.Error()})
					c.Writer.Flush()
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
			}
		}(driver.(db.Database))

		for {
			select {
			case msg, ok := <-client.StatusCode:

				if !ok {
					return false
				}
				c.SSEvent("message", map[string]interface{}{"current": msg, "total": len(stars)})
				c.Writer.Flush()
			case <-c.Request.Context().Done():
				fmt.Println("client closed")
				return false
			}
		}
	})
}
