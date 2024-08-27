package controller

import (
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/db"
	"github.com/weichen-lin/kafka-service/util"
)

var testDB *db.Database
var testController *Controller
var testJWTSecretKey = "secretfor32stringsecretfor32stringsecretfor32stringsecretfor32stringsecretfor32stringsecretfor32string"
var testJWTMaker util.Maker

func NewTestDatabase() *db.Database {
	driver, err := neo4j.NewDriverWithContext(
		"neo4j://localhost:7687",
		neo4j.BasicAuth("neo4j", "password", ""),
	)

	if err != nil {
		panic(err)
	}

	return &db.Database{
		Driver: driver,
	}
}

func NewTestJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		payload, err := testJWTMaker.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("email", payload.Email)

		c.Next()
	}
}

func TestMain(m *testing.M) {
	testDB = NewTestDatabase()

	testController = &Controller{
		db: testDB,
	}

	var err error
	testJWTMaker, err = util.NewJWTMaker(testJWTSecretKey)
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
