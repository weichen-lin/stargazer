package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/kafka-service/domain"
)

func Test_GetCrontabNoContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/crontab", testController.GetCrontab)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/crontab", nil)

	r.ServeHTTP(w, req)

	require.Equal(t, w.Code, http.StatusUnauthorized)
}

func Test_GetCrontab(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/crontab", NewTestJWTAuth(), testController.GetCrontab)

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

func Test_CreateCrontab(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/crontab", NewTestJWTAuth(), testController.CreateCrontab)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/crontab", nil)

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
	req, _ = http.NewRequest("POST", "/crontab", nil)
	req.Header.Set("Authorization", token)

	r.ServeHTTP(w, req)

	var response domain.CrontabEntity
	err = json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	validTime := time.Now().Format(time.RFC3339)

	require.Equal(t, "", response.TriggeredAt)
	require.Equal(t, "", response.UpdatedAt)
	require.Equal(t, validTime, response.CreatedAt)
	require.Equal(t, int64(1), response.Version)
	require.Equal(t, "new", response.Status)
	require.Equal(t, "", response.LastTriggeredAt)

	require.Equal(t, http.StatusCreated, w.Code)
}
