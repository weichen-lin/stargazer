package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/util"
)

func (c *Controller) SyncRepository(ctx *gin.Context) {
	user, err := c.db.GetUser(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	}

	limiter := fmt.Sprintf("%s-sync-limiter", user.Name())

	_, err = c.rdb.Get(ctx, limiter).Result()
	if err != nil {
		if err == redis.Nil {
			c.rdb.SetNX(ctx, limiter, true, time.Minute*30)
			c.kabaka.Publish("star-syncer", []byte(`{"email":"`+user.Email()+`","page":1}`), nil)
			ctx.JSON(http.StatusOK, "ok")
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server is not response now"})
		return
	}

	ttl, err := c.rdb.TTL(ctx, limiter).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server is not response now"})
	} else {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"expires": fmt.Sprintf("%d minutes", int(ttl.Minutes())),
			"message": "This user is already being processed. Please try again later.",
		})
	}
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

	results, err := c.db.SearchRepositoryByLanguageWithCollection(ctx, &db.SearchParams{
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

	topicCache := fmt.Sprintf("%s-topic-cache", user.Name())

	val, err := c.rdb.Get(ctx, topicCache).Result()
	if err != nil {
		if err == redis.Nil {
			results, err := c.db.GetAllRepositoryTopics(ctx)
			if err != nil {
				handleRepositoryErr(err, ctx)
				return
			}

			topicsMap := make(map[string][]int64)

			for _, result := range results {
				for _, topic := range result.Topics {
					repos, exists := topicsMap[topic]

					if !exists {
						repos := []int64{}
						repos = append(repos, result.RepoId)
						topicsMap[topic] = repos
						continue
					}

					repos = append(repos, result.RepoId)
					topicsMap[topic] = repos
				}
			}

			topics := util.GetRepositoryTopics(topicsMap)

			jsonString, _ := json.Marshal(topics)

			ctx.JSON(http.StatusOK, gin.H{
				"data": topics,
			})

			c.rdb.SetNX(ctx, topicCache, jsonString, time.Hour*12)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server is not response now"})
		return
	}

	var topics []*util.TopicsResult
	err = json.Unmarshal([]byte(val), &topics)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server is not response now"})
		return
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
