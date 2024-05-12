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

func (c *Controller) HandleConnections(ctx *gin.Context) {

	email, ok := ctx.Value("email").(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	stars, err := c.db.GetUserNotVectorize(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user stars"})
		return
	}

	client := &Client{
		Email:      email,
		StatusCode: make(chan int),
	}

	clients[email] = client
	defer delete(clients, email)

	ctx.Stream(func(w io.Writer) bool {
		w.(http.ResponseWriter).Header().Set("Content-Type", "text/event-stream")
		w.(http.ResponseWriter).Header().Set("Cache-Control", "no-cache")
		w.(http.ResponseWriter).Header().Set("Connection", "keep-alive")
		w.(http.ResponseWriter).Header().Set("Access-Control-Allow-Origin", "*")

		ctx.SSEvent("message", map[string]interface{}{"total": len(stars)})
		ctx.Writer.Flush()

		go func(db *db.Database) {
			defer close(client.StatusCode)
			for i := 0; i < len(stars); i++ {

				id := stars[i]

				_, err := workflow.VectorizeStar(&workflow.SyncUserStar{
					Email:  email,
					RepoId: id,
				})

				if err != nil {
					ctx.SSEvent("message", map[string]interface{}{"error": err.Error()})
					ctx.Writer.Flush()
					ctx.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				client.StatusCode <- i

				err = db.ConfirmVectorize(&workflow.SyncUserStar{
					Email:  email,
					RepoId: id,
				})

				if err != nil {
					ctx.SSEvent("message", map[string]interface{}{"error": err.Error()})
					ctx.Writer.Flush()
					ctx.AbortWithStatus(http.StatusInternalServerError)
					return
				}
			}
		}(c.db)

		for {
			select {
			case msg, ok := <-client.StatusCode:

				if !ok {
					return false
				}
				ctx.SSEvent("message", map[string]interface{}{"current": msg, "total": len(stars)})
				ctx.Writer.Flush()
			case <-ctx.Request.Context().Done():
				fmt.Println("client closed")
				return false
			}
		}
	})
}
