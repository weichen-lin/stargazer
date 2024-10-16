package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/stargazer/domain"
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
		user, token := createUserWithToken(t)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/crontab", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)

		var response domain.CrontabEntity
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		validTime := time.Now().Format(time.RFC3339)

		require.Equal(t, "", response.TriggeredAt)
		require.Equal(t, validTime, response.UpdatedAt)
		require.Equal(t, validTime, response.CreatedAt)
		require.Equal(t, "new", response.Status)
		require.Equal(t, "", response.LastTriggeredAt)

		require.Equal(t, http.StatusCreated, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/crontab", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)

		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		id := testController.scheduler.GetJob(user.Email())
		require.Equal(t, uuid.Nil, id)

		require.Equal(t, "", response.TriggeredAt)
		require.Equal(t, validTime, response.CreatedAt)
		require.Equal(t, validTime, response.UpdatedAt)
		require.Equal(t, "new", response.Status)
		require.Equal(t, "", response.LastTriggeredAt)

		require.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()

		newTriggeredAt := time.Now().UTC()
		req, _ = http.NewRequest("PATCH", fmt.Sprintf("/crontab?triggered_at=%s", newTriggeredAt.Format(time.RFC3339)), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		id = testController.scheduler.GetJob(user.Email())
		require.NotEqual(t, uuid.Nil, id)

		job := testController.scheduler.GetJobInfo(id)
		require.NotNil(t, job)

		nextRun, err := job.NextRun()
		require.NoError(t, err)

		location := newTriggeredAt.Location()
		nextRun = nextRun.In(location)

		require.Equal(t, newTriggeredAt.Hour(), nextRun.Hour())

		updatedTime := time.Now()

		checkUpdatedTime, err := time.Parse(time.RFC3339, response.UpdatedAt)
		require.NoError(t, err)

		require.Equal(t, newTriggeredAt.Format(time.RFC3339), response.TriggeredAt)
		require.WithinDuration(t, updatedTime, checkUpdatedTime, time.Second*2)
		require.Equal(t, validTime, response.CreatedAt)
		require.Equal(t, "new", response.Status)
		require.Equal(t, "", response.LastTriggeredAt)

		require.Equal(t, w.Code, http.StatusOK)
	})
}
