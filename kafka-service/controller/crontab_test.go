package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/kafka-service/domain"
)


func Test_GetCrontab(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/crontab", NewTestJWTAuth(), testController.GetCronTab)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/crontab", nil)

	r.ServeHTTP(w, req)

	require.Equal(t, w.Code, http.StatusUnauthorized)

	entity := &domain.UserEntity{
		Name:              faker.Name(),
		Email:             faker.Email(),
		Image:             "https://example.com/avatar.jpg",
		AccessToken:       "abc123",
		Provider:          "github",
		ProviderAccountId: "123456",
		Scope:             "read:user,user:email",
		AuthType:          "oauth",
		TokenType:         "bearer",
	}

	user := domain.FromUserEntity(entity)

	err := testDB.CreateUser(user)
	require.NoError(t, err)

	token, err := testJWTMaker.CreateToken(user.Email(), time.Now().Add(time.Hour))
	require.NoError(t, err)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/crontab", nil)
	req.Header.Set("Authorization", token)

	r.ServeHTTP(w, req)

	require.Equal(t, w.Code, http.StatusNotFound)
}
