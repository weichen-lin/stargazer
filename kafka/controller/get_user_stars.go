package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

type GetGithubReposInfo struct {
	UserId   string `form:"user_id" json:"user_id" binding:"required"`
	Username string `form:"user_name" json:"user_name" binding:"required"`
	Page     int    `form:"page" json:"page" binding:"required"`
}

var getUserStarsLimiter = cache.New(20*time.Minute, 10*time.Minute)

func GetUserStars(c *gin.Context) {
	var info GetGithubReposInfo

	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, found := getUserStarsLimiter.Get(info.UserId); found {

		_, expired, _ := getUserStarsLimiter.GetWithExpiration(info.UserId)
		remain := time.Until(expired)
		mins := int(remain.Minutes())

		c.JSON(http.StatusConflict, gin.H{
			"message": "This user is already being processed. Please try again later.",
			"expires": fmt.Sprintf("%d minutes", mins),
		})
		return
	}

	getUserStarsLimiter.Set(info.UserId, true, time.Minute*30)

	jsonString, err := json.Marshal(info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	producer, exists := c.Get("kafkaProducer")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kafka producer not found"})
		return
	}

	_, _, err = producer.(sarama.SyncProducer).SendMessage(&sarama.ProducerMessage{
		Topic: "get_user_stars",
		Value: sarama.StringEncoder(jsonString),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, gin.H{
		"message": "OK",
	})
}
