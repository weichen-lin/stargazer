package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/kafka-service/util"
)

type Middleware interface {
	BasicAuth() gin.HandlerFunc
	JWTAuth() gin.HandlerFunc
	Cors() gin.HandlerFunc
}

type middleware struct {
	secret string
	maker  util.Maker
}

func NewMiddleware() Middleware {
	secret := os.Getenv("AUTHENTICATION_TOKEN")
	if secret == "" {
		panic("AUTHENTICATION_TOKEN is not set")
	}

	jwtMaker, err := util.NewJWTMaker(secret)

	if err != nil {
		panic(err)
	}

	return &middleware{
		secret: secret,
		maker:  jwtMaker,
	}
}

func (m *middleware) BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if token != "Bearer "+m.secret {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *middleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		payload, err := m.maker.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("email", payload.Email)

		c.Next()
	}
}

func (m *middleware) Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
