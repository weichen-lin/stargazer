package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/kafka-service/domain"
)



func (c *Controller) GetRepository(ctx *gin.Context) {
	repo, err := c.db.GetRepository(ctx, 123456)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() })
		return
	}

	ctx.JSON(http.StatusOK, repo.ToRepositoryEntity())
}

func (c *Controller) CreateRepository(ctx *gin.Context) {

	createdAt := "2014-03-24T16:04:04Z"
	updatedAt := "2014-03-24T17:04:04Z"

	githubRepo := &domain.GithubRepository{
		ID:              1234567,
		Name:            "sample-repo",
		FullName:        "user/sample-repo",
		Owner:           domain.Owner{Login: "user", AvatarURL: "https://avatar.url", HTMLURL: "https://github.com/user"},
		HTMLURL:         "https://github.com/user/sample-repo",
		Description:     "This is a sample repository",
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		Homepage:        "https://sample-repo.com",
		StargazersCount: 10,
		WatchersCount:   20,
		Language:        "Go",
		Archived:        false,
		OpenIssuesCount: 5,
		Topics:          []string{"go", "sample"},
		DefaultBranch:   "main",
	}

	repo, _ := domain.NewRepository(githubRepo)

	err := c.db.SaveRepository(ctx, repo)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, repo.ToRepositoryEntity()) 
	}

	ctx.JSON(http.StatusOK, repo.ToRepositoryEntity())
}