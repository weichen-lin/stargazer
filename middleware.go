package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ClerkAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		sessionToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		// Verify the session
		claims, err := jwt.Verify(c.Request.Context(), &jwt.VerifyParams{
			Token: sessionToken,
		})
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		fmt.Println(claims.Subject)

		c.Set("clerkId", claims.Subject)
		// c.Set("userId", "user_2uiCzCED6xtVkHyQI1YzGU1aFaO")
	}
}

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		traceID := c.Request.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		c.Set("traceId", traceID)
		c.Writer.Header().Set("X-Trace-ID", traceID)

		c.Next()
	}
}

func BackgroundAuthMiddleware(bgToken string) gin.HandlerFunc {

	return func(c *gin.Context) {
		token := c.GetHeader("X-StarGazer-Token")
		if token != bgToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func ContextProvider(userId string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userId", userId)
	}
}
