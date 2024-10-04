package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/domain"
)

type GetCollectionQuery struct {
	Id string `form:"id" binding:"required"`
}

func (c *Controller) GetCollection(ctx *gin.Context) {
	id := ctx.Param("id")

	email, ok := db.EmailFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sharedCollection, err := c.db.GetCollectionById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "collection not found"})
		return
	}

	if sharedCollection.Owner != email && sharedCollection.SharedFrom == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx.JSON(http.StatusOK, sharedCollection)
}

type CreateCollectionRequest struct {
	Name string `json:"name" binding:"required"`
}

func (c *Controller) CreateCollection(ctx *gin.Context) {
	var body CreateCollectionRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	checkCollection, _ := c.db.GetCollectionByName(ctx, body.Name)
	if checkCollection != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "collection already exists"})
		return
	}

	collection, _ := domain.NewCollection(body.Name)

	err := c.db.SaveCollection(ctx, collection)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, collection.ToCollectionEntity())
}

type DeleteCollectionRequest struct {
	Id string `json:"id" binding:"required"`
}

func (c *Controller) DeleteCollection(ctx *gin.Context) {
	var body DeleteCollectionRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sharedCollection, err := c.db.GetCollectionById(ctx, body.Id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "collection not found"})
		return
	}

	err = c.db.DeleteCollection(ctx, sharedCollection.Collection.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "collection deleted"})
}

type CollectionRepoRequest struct {
	Id      string  `json:"id" binding:"required"`
	RepoIds []int64 `json:"repo_ids" binding:"required"`
}

func (c *Controller) AddRepoIntoCollection(ctx *gin.Context) {
	var body CollectionRepoRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sharedCollection, err := c.db.GetCollectionById(ctx, body.Id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "collection not found"})
		return
	}

	collection, _ := domain.FromCollectionEntity(sharedCollection.Collection)

	err = c.db.AddRepoToCollection(ctx, collection, body.RepoIds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, collection.ToCollectionEntity())
}

func (c *Controller) RemoveRepoFromCollection(ctx *gin.Context) {
	var body CollectionRepoRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sharedCollection, err := c.db.GetCollectionById(ctx, body.Id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "collection not found"})
		return
	}

	collection, _ := domain.FromCollectionEntity(sharedCollection.Collection)

	err = c.db.DeleteRepoFromCollection(ctx, collection, body.RepoIds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, collection.ToCollectionEntity())
}

type SearchRepoAtCollectionQuery struct {
	Page  int64  `form:"page" binding:"required"`
	Limit int64  `form:"limit" binding:"required"`
	Id    string `form:"id" binding:"required"`
}

func (c Controller) GetReposInCollection(ctx *gin.Context) {
	var query SearchRepoAtCollectionQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sharedCollection, err := c.db.GetCollectionById(ctx, query.Id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "collection not found"})
		return
	}

	collection, _ := domain.FromCollectionEntity(sharedCollection.Collection)

	repos, err := c.db.GetCollectionContainRepos(ctx, collection, query.Page, query.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, repos)
}

type SearchCollectionQuery struct {
	Page  int64  `form:"page" binding:"required"`
	Limit int64  `form:"limit" binding:"required"`
}

func (c Controller) GetCollections(ctx *gin.Context) {
	var query SearchCollectionQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.db.GetCollections(ctx, &db.PagingParams{
		Page:  query.Page,
		Limit: query.Limit,
	})

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "collection not found"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
