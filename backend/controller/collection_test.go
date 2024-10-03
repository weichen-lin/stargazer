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
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/domain"
)

func Test_GetCollection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/collection", NewTestJWTAuth(), testController.GetCollections)
	r.GET("/collection/:id", NewTestJWTAuth(), testController.GetCollection)

	user, token := createUserWithToken(t)

	t.Run("Unauthorized request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/collection", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test invalid query", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/collection?invalid=asda", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/collection/asdasd", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test real result", func(t *testing.T) {
		collection := createCollection(t, user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/collection/%s", uuid.New().String()), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)

		var response *domain.CollectionEntity
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/collection/%s", collection.Id().String()), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, collection.Id().String(), response.Id)
		require.Equal(t, collection.Name(), response.Name)
	})
}

func Test_CreateAndRemoveCollection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/collection", NewTestJWTAuth(), testController.GetCollections)
	r.GET("/collection/:id", NewTestJWTAuth(), testController.GetCollection)
	r.POST("/collection", NewTestJWTAuth(), testController.CreateCollection)
	r.DELETE("/collection", NewTestJWTAuth(), testController.DeleteCollection)

	user, token := createUserWithToken(t)

	t.Run("Unauthorized request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/collection", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/collection", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test invalid body", func(t *testing.T) {
		body, err := json.Marshal(`{"invalid": "value"}`)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/collection", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/collection", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/collection?page=asda&limit=1", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/collection?page=1&limit=asd", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test real result", func(t *testing.T) {
		collection := createCollection(t, user)

		duplicateBody, err := json.Marshal(&CreateCollectionRequest{
			Name: collection.Name(),
		})
		require.NoError(t, err)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/collection", bytes.NewBuffer(duplicateBody))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusConflict, w.Code)

		test_name := faker.Name()
		body, err := json.Marshal(&CreateCollectionRequest{
			Name: test_name,
		})
		require.NoError(t, err)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/collection", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)

		var response *domain.CollectionEntity
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/collection/%s", response.Id), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, response.Name, test_name)

		body, err = json.Marshal(&DeleteCollectionRequest{
			Id: response.Id,
		})
		require.NoError(t, err)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/collection", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/collection/%s", response.Id), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Test get collections", func(t *testing.T) {
		user, token := createUserWithToken(t)

		for i := 0; i < 25; i++ {
			createCollection(t, user)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/collection?page=1&limit=20", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var response *db.CollectionSearchResult
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		require.Len(t, response.Data, 20)
		require.Equal(t, int64(25), response.Total)
	})
}

func Test_AddAndRemoveRepoIntoCollection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/collection/repos", NewTestJWTAuth(), testController.GetReposInCollection)
	r.POST("/collection/repos", NewTestJWTAuth(), testController.AddRepoIntoCollection)
	r.DELETE("/collection/repos", NewTestJWTAuth(), testController.RemoveRepoFromCollection)

	user, token := createUserWithToken(t)

	t.Run("Unauthorized request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/collection/repos", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/collection/repos", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test invalid body", func(t *testing.T) {
		body, err := json.Marshal(`{"repo_ids": ["value"]}`)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/collection/repos", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/collection/repos", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test real result", func(t *testing.T) {
		collection := createCollection(t, user)

		repos := make([]int64, 25)
		for i := 0; i < 25; i++ {
			repo := createRepository(t, user)
			repos[i] = repo.RepoID()
		}

		body, err := json.Marshal(&CollectionRepoRequest{
			Id:      collection.Id().String(),
			RepoIds: repos,
		})
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/collection/repos", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/collection/repos?id=%s&page=1&limit=20", collection.Id().String()), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var response *db.SearchResult
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		require.Len(t, response.Data, 20)
		require.Equal(t, response.Total, int64(25))

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/collection/repos", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/collection/repos?id=%s&page=1&limit=20", collection.Id().String()), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		require.Len(t, response.Data, 0)
		require.Equal(t, response.Total, int64(0))
	})
}
