package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func Neo4jDriverMiddleware(driver neo4j.DriverWithContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("neo4jDriver", driver)
		c.Next()
	}
}
