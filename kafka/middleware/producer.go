package middleware

import (
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

func ProducerMiddleware(producer sarama.SyncProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("kafkaProducer", producer)
		c.Next()
	}
}
