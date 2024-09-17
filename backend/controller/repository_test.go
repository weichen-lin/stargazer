package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/domain"
	"golang.org/x/exp/rand"
)

var Languages []string = []string{
	"Go",
	"Python",
	"JavaScript",
	"Java",
	"C++",
	"C#",
	"Ruby",
	"Swift",
	"Kotlin",
	"TypeScript",
	"Rust",
	"PHP",
	"Scala",
	"Haskell",
	"Perl",
	"R",
	"MATLAB",
	"Dart",
	"Lua",
	"Groovy",
}

func getRandomLanguage() string {
	seed := uint64(time.Now().UnixNano()) & ((1 << 63) - 1)
	rand.Seed(seed)

	randomIndex := rand.Intn(len(Languages))

	return Languages[randomIndex]
}

func createRepository(t *testing.T, user *domain.User) *domain.Repository {
	repositoryEntity := &domain.RepositoryEntity{
		RepoID:            faker.RandomUnixTime(),
		RepoName:          faker.Name(),
		OwnerName:         faker.Name(),
		AvatarURL:         faker.URL(),
		HtmlURL:           faker.URL(),
		Homepage:          faker.URL(),
		CreatedAt:         "2024-01-01T00:00:00Z",
		UpdatedAt:         "2024-01-02T00:00:00Z",
		StargazersCount:   100,
		WatchersCount:     50,
		OpenIssuesCount:   10,
		DefaultBranch:     "main",
		Description:       faker.Sentence(),
		Language:          getRandomLanguage(),
		Archived:          false,
		ExternalCreatedAt: time.Now().Format(time.RFC3339),
		LastSyncedAt:      time.Now().Format(time.RFC3339),
		LastModifiedAt:    time.Now().Format(time.RFC3339),
	}

	repo, err := domain.FromRepositoryEntity(repositoryEntity)
	require.NoError(t, err)
	require.NotEmpty(t, repo)

	ctx, err := db.WithEmail(context.Background(), user.Email())
	require.NoError(t, err)
	require.NotEmpty(t, ctx)

	err = testDB.CreateRepository(ctx, repo)
	require.NoError(t, err)

	return repo
}

func Test_GetRepository(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/repository/:id", NewTestJWTAuth(), testController.GetRepository)

	user, token := createUserWithToken(t)

	t.Run("Unauthorized request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository/123123", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test not exist repo_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository/123123", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Test invalid repo_id format", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository/asdasdasd", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test get repo success", func(t *testing.T) {
		repo := createRepository(t, user)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/repository/%d", repo.RepoID()), nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})
}

func Test_GetRepositoryLanguageDistribution(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/repository/language-distribution", NewTestJWTAuth(), testController.GetUserLanguageDistribution)

	user, token := createUserWithToken(t)

	t.Run("Unauthorized request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository/language-distribution", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test empty", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository/language-distribution", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var response []db.LanguageDistribution
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, 0, len(response))
	})

	t.Run("Test real language distribution", func(t *testing.T) {
		recordMap := make(map[string]int)

		for i := 0; i < 100; i++ {
			repo := createRepository(t, user)

			if count, exists := recordMap[repo.Language()]; exists {
				recordMap[repo.Language()] = count + 1
			} else {
				recordMap[repo.Language()] = 1
			}
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository/language-distribution", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var response []db.LanguageDistribution
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		for _, data := range response {
			count, exists := recordMap[data.Language]
			require.Equal(t, true, exists)
			require.Equal(t, int64(count), data.Count)
		}
	})
}

func contains(arr []string, language string) bool {
	for _, lang := range arr {
		if language == lang {
			return true
		}
	}

	return false
}

func Test_GetRepositoryByPageAndLanguages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/repository", NewTestJWTAuth(), testController.SearchRepoByLanguages)

	user, token := createUserWithToken(t)

	t.Run("Unauthorized request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test invalid query", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository?page=1aasd&limit=aa&languages=123,123123,123", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/repository?page=1&limit=aa&languages=123,123123,123", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test real result", func(t *testing.T) {
		recordMap := make(map[string]int)

		for i := 0; i < 100; i++ {
			repo := createRepository(t, user)

			if count, exists := recordMap[repo.Language()]; exists {
				recordMap[repo.Language()] = count + 1
			} else {
				recordMap[repo.Language()] = 1
			}
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository?page=1&limit=20&languages=Go,Python", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var response db.SearchResult
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		realTotal := recordMap["Go"] + recordMap["Python"]
		require.Equal(t, response.Total, int64(realTotal))

		for _, data := range response.Data {
			require.Equal(t, true, contains([]string{"Go", "Python"}, data.Language))
		}
	})
}

func Test_GetRepositoryBySortParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/repository/sort", NewTestJWTAuth(), testController.GetRepositoriesByKey)

	user, token := createUserWithToken(t)

	t.Run("Unauthorized request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository/sort", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test invalid query", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository/sort?key=1aasd&order=aa", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/repository/sort?key=created_at&order=aa", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test real result", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			repo := createRepository(t, user)
			require.NotEmpty(t, repo)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository/sort?key=created_at&order=DESC", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data []*domain.Repository `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		require.Equal(t, len(response.Data), 5)
	})
}

func Test_FullTextSearchByQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/repository/full-text-search", NewTestJWTAuth(), testController.FullTextSearchWithQuery)

	user, token := createUserWithToken(t)

	t.Run("Unauthorized request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository/full-text-search", nil)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Test empty query", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository/full-text-search?query=", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test real result", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			repo := createRepository(t, user)
			require.NotEmpty(t, repo)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/repository/full-text-search?query=et", nil)
		req.Header.Set("Authorization", token)

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var response []*domain.Repository
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
	})
}
