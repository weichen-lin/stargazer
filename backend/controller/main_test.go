package controller

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/domain"
	"github.com/weichen-lin/stargazer/util"
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

func createUserWithToken(t *testing.T) (*domain.User, string) {
	entity := &domain.UserEntity{
		Name:              faker.Name(),
		Email:             faker.Email(),
		Image:             faker.URL(),
		AccessToken:       faker.Sentence(),
		Provider:          faker.Sentence(),
		ProviderAccountId: faker.Sentence(),
		Scope:             faker.Sentence(),
		AuthType:          faker.Sentence(),
		TokenType:         faker.Sentence(),
	}

	user := domain.FromUserEntity(entity)

	err := testDB.CreateUser(user)
	require.NoError(t, err)

	token, err := testJWTMaker.CreateToken(user.Email(), time.Now().Add(time.Hour))
	require.NoError(t, err)

	return user, token
}

func createFolder(t *testing.T, user *domain.User) *domain.Folder {
	folder, err := domain.NewFolder(faker.Name())
	require.NotEmpty(t, folder)

	ctx, err := db.WithEmail(context.Background(), user.Email())
	require.NoError(t, err)
	require.NotEmpty(t, ctx)

	err = testDB.SaveFolder(ctx, folder)
	require.NoError(t, err)

	return folder
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
