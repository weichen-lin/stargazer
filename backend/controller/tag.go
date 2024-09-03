package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/stargazer/domain"
)

func (c *Controller) GetTags(ctx *gin.Context) {
	id := ctx.Param("id")

	repo_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid repo_id"})
		return
	}

	tags, _ := c.db.GetTagsByRepo(ctx, repo_id)

	tagEntities := make([]*domain.TagEntity, len(tags))

	for i, tag := range tags {
		tagEntities[i] = tag.ToTagEntity()
	}

	ctx.JSON(http.StatusOK, tagEntities)
}

type TagRequest struct {
	Name   string `json:"name" binding:"required"`
	RepoId int64  `json:"repo_id" binding:"required,min=1"`
}

func (c *Controller) CreateTag(ctx *gin.Context) {
	var body TagRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := domain.NewTag(body.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.db.SaveTag(ctx, tag, body.RepoId)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, tag.ToTagEntity())
}

func (c *Controller) DeleteTag(ctx *gin.Context) {
	var body TagRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := c.db.GetTagByName(ctx, body.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.db.RemoveTag(ctx, tag, body.RepoId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "tag deleted"})
}
