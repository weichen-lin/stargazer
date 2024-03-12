package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
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

	var job SyncUserStars
	if err := c.ShouldBind(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := &Client{
		UserName:   job.Username,
		StatusCode: make(chan int, len(job.Stars)),
	}

	clients[job.Username] = client
	defer delete(clients, job.Username)

	c.Stream(func(w io.Writer) bool {
		w.(http.ResponseWriter).Header().Set("Content-Type", "text/event-stream")
		w.(http.ResponseWriter).Header().Set("Cache-Control", "no-cache")
		w.(http.ResponseWriter).Header().Set("Connection", "keep-alive")
		w.(http.ResponseWriter).Header().Set("Access-Control-Allow-Origin", "*")

		go func() {
			defer close(client.StatusCode)
			for i := 0; i < len(job.Stars); i++ {

				id := job.Stars[i]

				status, _ := workflow.VectorizeStar(&workflow.SyncUserStarMsg{
					UserName: job.Username,
					RepoId:   id,
				})

				select {
				case client.StatusCode <- status:
                    fmt.Println("vectorize star status: ", status)
				case <-c.Request.Context().Done():
					return
				}
			}
		}()

		for {
			select {
			case msg, ok := <-client.StatusCode:
				if !ok {
					return false
				}
				c.SSEvent("message", gin.H{"data": msg})
			case <-c.Request.Context().Done():
				return false
			}
		}
	})
}
