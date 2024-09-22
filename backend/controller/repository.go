package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/util"
)

var getUserStarsLimiter = cache.New(20*time.Minute, 10*time.Minute)
var getTopicsLimiter = cache.New(5*time.Minute, 5*time.Minute)

func (c *Controller) SyncRepository(ctx *gin.Context) {
	user, err := c.db.GetUser(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	}

	if _, found := getUserStarsLimiter.Get(user.Email()); found {
		_, expired, _ := getUserStarsLimiter.GetWithExpiration(user.Email())
		remain := time.Until(expired)
		mins := int(remain.Minutes())

		ctx.JSON(http.StatusConflict, gin.H{
			"message": "This user is already being processed. Please try again later.",
			"expires": fmt.Sprintf("%d minutes", mins),
		})
		return
	}

	getUserStarsLimiter.Set(user.Email(), true, time.Minute*30)

	c.kabaka.Publish("star-syncer", []byte(`{"email":"`+user.Email()+`","page":1}`))

	ctx.JSON(http.StatusOK, "ok")
}

type GetRepositoryQuery struct {
	RepoId int64 `form:"repo_id" binding:"required"`
}

func handleRepositoryErr(err error, ctx *gin.Context) {
	switch err {
	case db.ErrNotFoundEmailAtContext:
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UnAuthorized"})
	case db.ErrRepositoryNotFound:
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case db.ErrInvalidSortKey, db.ErrInvalidSortOrder:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func (c *Controller) DeleteRepository(ctx *gin.Context) {
	id := ctx.Param("id")

	repo_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid repo_id"})
		return
	}

	err = c.db.DeleteRepository(ctx, repo_id)

	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
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

func (c *Controller) GetTopics(ctx *gin.Context) {
	user, err := c.db.GetUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	}

	topics, _ := util.StarGazerTopicCache.GetTopics(user.Email())

	if _, found := getTopicsLimiter.Get(user.Email()); !found {
		c.kabaka.Publish("topic-syncer", []byte(`{"email":"`+user.Email()+`"}`))
		getTopicsLimiter.Set(user.Email(), true, time.Minute*5)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": topics,
	})
}

type GetRepositoriesByKeyQueries struct {
	Key   string `form:"key" binding:"required"`
	Order string `form:"order" binding:"required"`
}

func (c *Controller) GetRepositoriesByKey(ctx *gin.Context) {
	var queries GetRepositoriesByKeyQueries
	if err := ctx.ShouldBindQuery(&queries); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.db.GetRepositoriesOrderBy(ctx, &db.SortParams{
		Key:   queries.Key,
		Order: queries.Order,
	})

	if err != nil {
		fmt.Println(err)
		handleRepositoryErr(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

type FullTextSearchQuery struct {
	Query string `form:"query" binding:"required"`
}

func (c *Controller) FullTextSearchWithQuery(ctx *gin.Context) {
	var q FullTextSearchQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if q.Query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "not contains query string"})
		return
	}

	repos, err := c.db.FullTextSearch(ctx, q.Query)

	if err != nil {
		handleRepositoryErr(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, repos)
}
