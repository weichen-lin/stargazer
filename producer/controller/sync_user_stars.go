package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jOpeartion "github.com/weichen-lin/kafka-service/neo4j"
	"github.com/weichen-lin/kafka-service/workflow"
)

type Client struct {
	UserName   string
	StatusCode chan int
}

type SyncUserStars struct {
	Stars    []int  `json:"stars" binding:"required" form:"stars"`
	Username string `json:"username" binding:"required" form:"username"`
}

var clients = make(map[string]*Client)

func HandleConnections(c *gin.Context) {

	userName, ok := c.Value("userName").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	neo4jDriver, exists := c.Get("neo4jDriver")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neo4j driver not found"})
		return
	}

	_, ok = neo4jDriver.(neo4j.DriverWithContext)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Neo4j driver"})
		return
	}

	stars, err := neo4jOpeartion.GetUserNotVectorize(neo4jDriver.(neo4j.DriverWithContext), userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user stars"})
		return
	}

	client := &Client{
		UserName:   userName,
		StatusCode: make(chan int),
	}

	clients[userName] = client
	defer delete(clients, userName)

	c.Stream(func(w io.Writer) bool {
		w.(http.ResponseWriter).Header().Set("Content-Type", "text/event-stream")
		w.(http.ResponseWriter).Header().Set("Cache-Control", "no-cache")
		w.(http.ResponseWriter).Header().Set("Connection", "keep-alive")
		w.(http.ResponseWriter).Header().Set("Access-Control-Allow-Origin", "*")

		c.SSEvent("message", map[string]interface{}{"total": len(stars)})
		c.Writer.Flush()

		go func(driver neo4j.DriverWithContext) {
			defer close(client.StatusCode)
			for i := 0; i < len(stars); i++ {

				id := stars[i]

				_, err := workflow.VectorizeStar(&workflow.SyncUserStarMsg{
					UserName: userName,
					RepoId:   id,
				})

				if err != nil {
					c.SSEvent("message", map[string]interface{}{"error": err.Error()})
					c.Writer.Flush()
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				client.StatusCode <- i

				err = neo4jOpeartion.ConfirmVectorize(driver, &workflow.SyncUserStarMsg{
					UserName: userName,
					RepoId:   id,
				})

				if err != nil {
					c.SSEvent("message", map[string]interface{}{"error": err.Error()})
					c.Writer.Flush()
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
			}
		}(neo4jDriver.(neo4j.DriverWithContext))

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
