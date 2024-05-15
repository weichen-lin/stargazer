package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

var getUserStarsLimiter = cache.New(20*time.Minute, 10*time.Minute)

func (c *Controller) GetUserStars(ctx *gin.Context) {

	email, ok := ctx.Value("email").(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if _, found := getUserStarsLimiter.Get(email); found {
		_, expired, _ := getUserStarsLimiter.GetWithExpiration(email)
		remain := time.Until(expired)
		mins := int(remain.Minutes())

		ctx.JSON(http.StatusConflict, gin.H{
			"message": "This user is already being processed. Please try again later.",
			"expires": fmt.Sprintf("%d minutes", mins),
		})
		return
	}

	getUserStarsLimiter.Set(email, true, time.Minute*30)

	_, _, err := c.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "get_user_stars",
		Value: sarama.StringEncoder(`{"email":"` + email + `","page":1}`),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(200, gin.H{
		"message": "OK",
	})
}
