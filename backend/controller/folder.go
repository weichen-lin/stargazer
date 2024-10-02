package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/stargazer/domain"
)

type GetFolderQuery struct {
	Id string `form:"id" binding:"required"`
}

func (c *Controller) GetFolder(ctx *gin.Context) {
	var query GetFolderQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder, err := c.db.GetFolderById(ctx, query.Id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}

	ctx.JSON(http.StatusOK, folder.ToFolderEntity())
}

type CreateFolderRequest struct {
	Name string `json:"name" binding:"required"`
}

func (c *Controller) CreateFolder(ctx *gin.Context) {
	var body CreateFolderRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	checkFolder, _ := c.db.GetFolderByName(ctx, body.Name)
	if checkFolder != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "folder already exists"})
		return
	}

	folder, _ := domain.NewFolder(body.Name)

	err := c.db.SaveFolder(ctx, folder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, folder.ToFolderEntity())
}

type DeleteFolderRequest struct {
	Id string `json:"id" binding:"required"`
}

func (c *Controller) DeleteFolder(ctx *gin.Context) {
	var body DeleteFolderRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder, err := c.db.GetFolderById(ctx, body.Id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}

	err = c.db.DeleteFolder(ctx, folder.Id().String())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "folder deleted"})
}

type FolderRepoRequest struct {
	Id      string  `json:"id" binding:"required"`
	RepoIds []int64 `json:"repo_ids" binding:"required"`
}

func (c *Controller) AddRepoIntoFolder(ctx *gin.Context) {
	var body FolderRepoRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder, err := c.db.GetFolderById(ctx, body.Id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}

	err = c.db.AddRepoToFolder(ctx, folder, body.RepoIds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, folder.ToFolderEntity())
}

func (c *Controller) RemoveRepoFromFolder(ctx *gin.Context) {
	var body FolderRepoRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder, err := c.db.GetFolderById(ctx, body.Id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}

	err = c.db.DeleteRepoFromFolder(ctx, folder, body.RepoIds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, folder.ToFolderEntity())
}

type SearchRepoAtFolderQuery struct {
	Page  int64  `form:"page" binding:"required"`
	Limit int64  `form:"limit" binding:"required"`
	Id    string `form:"id" binding:"required"`
}

func (c Controller) GetReposInFolder(ctx *gin.Context) {
	var query SearchRepoAtFolderQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder, err := c.db.GetFolderById(ctx, query.Id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}

	repos, err := c.db.GetFolderContainRepos(ctx, folder, query.Page, query.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, repos)
}
