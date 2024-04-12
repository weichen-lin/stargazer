package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/kafka-service/util"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		expectedToken := os.Getenv("AUTHENTICATION_TOKEN")

		if token != "Bearer "+expectedToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		expectedToken := os.Getenv("AUTHENTICATION_TOKEN")

		jwtMaker, err := util.NewJWTMaker(expectedToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		payload, err := jwtMaker.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("userName", payload.UserName)

		c.Next()
	}
}
