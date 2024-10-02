package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/domain"
)

func Test_GetFolder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/folder", NewTestJWTAuth(), testController.GetFolder)

	user, token := createUserWithToken(t)

	t.Run("Unauthorized request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/folder", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test invalid query", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/folder?invalid=asda", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/folder?id=", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test real result", func(t *testing.T) {
		folder := createFolder(t, user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/folder?id=%s", "invalid-id"), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)

		var response *domain.FolderEntity
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/folder?id=%s", folder.Id().String()), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, folder.Id().String(), response.Id)
		require.Equal(t, folder.Name(), response.Name)
	})
}

func Test_CreateAndRemoveFolder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/folder", NewTestJWTAuth(), testController.GetFolder)
	r.POST("/folder", NewTestJWTAuth(), testController.CreateFolder)
	r.DELETE("/folder", NewTestJWTAuth(), testController.DeleteFolder)

	user, token := createUserWithToken(t)

	t.Run("Unauthorized request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/folder", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/folder", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test invalid body", func(t *testing.T) {
		body, err := json.Marshal(`{"invalid": "value"}`)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/folder", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/folder?name=", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test real result", func(t *testing.T) {
		folder := createFolder(t, user)

		duplicateBody, err := json.Marshal(&CreateFolderRequest{
			Name: folder.Name(),
		})
		require.NoError(t, err)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/folder", bytes.NewBuffer(duplicateBody))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusConflict, w.Code)

		test_name := faker.Name()
		body, err := json.Marshal(&CreateFolderRequest{
			Name: test_name,
		})
		require.NoError(t, err)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/folder", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)

		var response *domain.FolderEntity
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/folder?id=%s", response.Id), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, response.Name, test_name)

		body, err = json.Marshal(&DeleteFolderRequest{
			Id: response.Id,
		})
		require.NoError(t, err)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/folder", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/folder?id=%s", response.Id), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
	})
}

func Test_AddAndRemoveRepoIntoFolder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/folder/repos", NewTestJWTAuth(), testController.GetReposInFolder)
	r.POST("/folder/repos", NewTestJWTAuth(), testController.AddRepoIntoFolder)
	r.DELETE("/folder/repos", NewTestJWTAuth(), testController.RemoveRepoFromFolder)

	user, token := createUserWithToken(t)

	t.Run("Unauthorized request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/folder/repos", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/folder/repos", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test invalid body", func(t *testing.T) {
		body, err := json.Marshal(`{"repo_ids": ["value"]}`)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/folder/repos", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/folder/repos", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test real result", func(t *testing.T) {
		folder := createFolder(t, user)

		repos := make([]int64, 25)
		for i := 0; i < 25; i++ {
			repo := createRepository(t, user)
			repos[i] = repo.RepoID()
		}

		body, err := json.Marshal(&FolderRepoRequest{
			Id:      folder.Id().String(),
			RepoIds: repos,
		})
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/folder/repos", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/folder/repos?id=%s&page=1&limit=20", folder.Id().String()), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var response *db.SearchResult
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		require.Len(t, response.Data, 20)
		require.Equal(t, response.Total, int64(25))

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/folder/repos", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/folder/repos?id=%s&page=1&limit=20", folder.Id().String()), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		require.Len(t, response.Data, 0)
		require.Equal(t, response.Total, int64(0))
	})
}
