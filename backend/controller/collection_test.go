package controller

import (
	"bytes"
	"context"
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

		var response *db.SharedCollection
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/collection/%s", collection.Id().String()), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, collection.Id().String(), response.Collection.Id)
		require.Equal(t, user.Email(), response.Owner)
		require.Nil(t, response.SharedFrom)
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

		var getResponse *db.SharedCollection
		err = json.NewDecoder(w.Body).Decode(&getResponse)
		require.NoError(t, err)
		require.Equal(t, response.Id, getResponse.Collection.Id)
		require.Equal(t, user.Email(), getResponse.Owner)
		require.Nil(t, getResponse.SharedFrom)

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

func Test_UpdateCollection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/collection/:id", NewTestJWTAuth(), testController.GetCollection)
	r.PATCH("/collection/:id", NewTestJWTAuth(), testController.UpdateCollection)

	user, token := createUserWithToken(t)
	collection := createCollection(t, user)

	ctx, err := db.WithEmail(context.Background(), user.Email())
	require.Equal(t, ctx.Value("email"), user.Email())
	require.NoError(t, err)

	t.Run("Unauthorized request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/collection/123123123", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test invalid body", func(t *testing.T) {
		invalidMessages := []string{
			`{"named": "test"}`,
			`{"nasme": "asdas", "dsescription": "asdasd", "is_public": true}`,
			`{"desscription": ""}`,
			`{"is_pusblic": "true"}`,
			`{"is_public": "true"}`,
		}

		for _, message := range invalidMessages {
			body := bytes.NewBufferString(message)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", fmt.Sprintf("/collection/%s", collection.Id().String()), body)
			req.Header.Set("Authorization", token)

			r.ServeHTTP(w, req)
			require.Equal(t, http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Test real result", func(t *testing.T) {
		validMessages := []string{
			`{"name": "new name"}`,
			`{"name": "new name", "description": "new description"}`,
			`{"name": "new name", "is_public": true}`,
			`{"name": "new name 2", "description": "new description 2", "is_public": false}`,
		}

		for _, message := range validMessages {
			var payload UpdateCollectionPayload
			err := json.Unmarshal([]byte(message), &payload)
			require.NoError(t, err)

			body := bytes.NewBufferString(message)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", fmt.Sprintf("/collection/%s", collection.Id()), body)
			req.Header.Set("Authorization", token)

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)

			var response *domain.CollectionEntity
			err = json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)

			require.Equal(t, collection.Id().String(), response.Id)
			require.Equal(t, payload.Name, response.Name)
			require.Equal(t, payload.Description, response.Description)
			require.Equal(t, payload.IsPublic, response.IsPublic)

			sharedCollection, err := testController.db.GetCollectionById(ctx, collection.Id().String())
			require.NoError(t, err)

			require.Equal(t, response.Name, sharedCollection.Collection.Name)
			require.Equal(t, response.Description, sharedCollection.Collection.Description)
			require.Equal(t, response.IsPublic, sharedCollection.Collection.IsPublic)
		}
	})
}
