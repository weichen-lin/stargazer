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
	"github.com/weichen-lin/stargazer/domain"
)

func Test_TagCRUD(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/tag/:id", NewTestJWTAuth(), testController.GetTags)
	r.POST("tag", NewTestJWTAuth(), testController.CreateTag)
	r.DELETE("/tag", NewTestJWTAuth(), testController.DeleteTag)

	user, token := createUserWithToken(t)
	repo := createRepository(t, user)

	t.Run("Unauthorized request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tag/123123", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/tag", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/tag", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid request body at create tags", func(t *testing.T) {
		tests := []struct {
			name    string
			payload interface{}
		}{
			{
				name: "Missing Name",
				payload: TagRequest{
					RepoId: 1,
				},
			},
			{
				name: "Missing RepoID",
				payload: TagRequest{
					Name: "TestTag",
				},
			},
			{
				name: "Invalid RepoID type",
				payload: map[string]interface{}{
					"Name":   "TestTag",
					"RepoID": "not a number",
				},
			},
			{
				name: "Negative RepoID",
				payload: TagRequest{
					Name:   "TestTag",
					RepoId: -1,
				},
			},
			{
				name: "Empty Name",
				payload: TagRequest{
					Name:   "",
					RepoId: 1,
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, err := json.Marshal(tt.payload)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/tag", bytes.NewBuffer(body))
				req.Header.Set("Authorization", token)
				req.Header.Set("Content-Type", "application/json")

				r.ServeHTTP(w, req)

				require.Equal(t, http.StatusBadRequest, w.Code, "Test case: %s", tt.name)
			})
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, err := json.Marshal(tt.payload)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("DELETE", "/tag", bytes.NewBuffer(body))
				req.Header.Set("Authorization", token)
				req.Header.Set("Content-Type", "application/json")

				r.ServeHTTP(w, req)

				require.Equal(t, http.StatusBadRequest, w.Code, "Test case: %s", tt.name)
			})
		}
	})

	t.Run("Get tags after create tag then delete it", func(t *testing.T) {
		payload := &TagRequest{
			Name:   faker.Name(),
			RepoId: repo.RepoID(),
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tag", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		var response domain.TagEntity
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, payload.Name, response.Name)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/tag/%d", repo.RepoID()), nil)
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)
		var getResponse []*domain.TagEntity
		err = json.NewDecoder(w.Body).Decode(&getResponse)
		require.NoError(t, err)
		require.Equal(t, 1, len(getResponse))

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/tag", bytes.NewBuffer(body))
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/tag/%d", repo.RepoID()), nil)
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)
		err = json.NewDecoder(w.Body).Decode(&getResponse)
		require.NoError(t, err)
		require.Equal(t, 0, len(getResponse))
	})
}
