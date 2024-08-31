package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/stargazer/db"
)

type GetRepositoryQuery struct {
	RepoId int64 `form:"repo_id" binding:"required"`
}

func handleRepositoryErr(err error, ctx *gin.Context) {
	switch err {
	case db.ErrNotFoundEmailAtContext:
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UnAuthorized"})
	case db.ErrRepositoryNotFound:
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server is not response now"})
	}
}

func (c *Controller) GetRepository(ctx *gin.Context) {
	id := ctx.Param("id")

	repo_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid repo_id"})
		return
	}

	repo, err := c.db.GetRepository(ctx, repo_id)

	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, repo.ToRepositoryEntity())
		return
	case errors.Is(err, db.ErrRepositoryNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("repository %d not found", repo_id)})
		return
	case errors.Is(err, db.ErrNotFoundEmailAtContext):
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	default:
		ctx.JSON(http.StatusForbidden, gin.H{"error": ""})
		return
	}
}

func (c *Controller) GetUserLanguageDistribution(ctx *gin.Context) {
	distribution, err := c.db.GetRepoLanguageDistribution(ctx)

	if err != nil {
		handleRepositoryErr(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, distribution)
}

type SearchRepositoryQuery struct {
	Page      int64  `form:"page" binding:"required"`
	Limit     int64  `form:"limit" binding:"required"`
	Languages string `form:"languages" binding:"required"`
}

func (c *Controller) SearchRepoByLanguages(ctx *gin.Context) {

	var query SearchRepositoryQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := c.db.SearchRepositoryByLanguage(ctx, &db.SearchParams{
		Page:      query.Page,
		Limit:     query.Limit,
		Languages: strings.Split(query.Languages, ","),
	})

	if err != nil {
		handleRepositoryErr(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, results)
}
