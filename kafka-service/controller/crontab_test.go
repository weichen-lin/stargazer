package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/kafka-service/domain"
)

func Test_CrontabCRUD(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/crontab", NewTestJWTAuth(), testController.GetCrontab)
	r.POST("/crontab", NewTestJWTAuth(), testController.CreateCrontab)
	r.PATCH("/crontab", NewTestJWTAuth(), testController.UpdateCrontab)

	t.Run("Unauthorized request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/crontab", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/crontab", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("PATCH", "/crontab", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test Get Crontab before create, raise not found error", func(t *testing.T) {
		_, token := createUserWithToken(t)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/crontab", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Test create crontab and then check crontab exist, final test update", func(t *testing.T) {
		_, token := createUserWithToken(t)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/crontab", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)

		var response domain.CrontabEntity
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		validTime := time.Now().Format(time.RFC3339)

		require.Equal(t, "", response.TriggeredAt)
		require.Equal(t, "", response.UpdatedAt)
		require.Equal(t, validTime, response.CreatedAt)
		require.Equal(t, int64(1), response.Version)
		require.Equal(t, "new", response.Status)
		require.Equal(t, "", response.LastTriggeredAt)

		require.Equal(t, http.StatusCreated, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/crontab", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)

		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, "", response.TriggeredAt)
		require.Equal(t, "", response.UpdatedAt)
		require.Equal(t, validTime, response.CreatedAt)
		require.Equal(t, int64(1), response.Version)
		require.Equal(t, "new", response.Status)
		require.Equal(t, "", response.LastTriggeredAt)

		require.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("PATCH", "/crontab?hour=12", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)

		require.Equal(t, w.Code, http.StatusOK)

		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		updatedTime := time.Now().Format(time.RFC3339)

		triggeredTime, _ := getTime(12)

		require.Equal(t, triggeredTime.Format(time.RFC3339), response.TriggeredAt)
		require.Equal(t, updatedTime, response.UpdatedAt)
		require.Equal(t, validTime, response.CreatedAt)
		require.Equal(t, int64(2), response.Version)
		require.Equal(t, "new", response.Status)
		require.Equal(t, "", response.LastTriggeredAt)

		require.Equal(t, w.Code, http.StatusOK)
	})
}
