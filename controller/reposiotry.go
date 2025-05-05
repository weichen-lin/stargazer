package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/stargazer/domain"
)

func (c Controller) GetLanguageDistribution(ctx *gin.Context) {
	languageDistribution, err := c.service.GetLanguageDistribution(ctx)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, languageDistribution)
}

func (c Controller) GetLatestUserUpdatedRepositories(ctx *gin.Context) {
	repositories, err := c.service.GetLatestUserUpdatedRepositories(ctx)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	entities := make([]*domain.RepositoryEntity, 0, len(repositories))
	for _, repository := range repositories {
		entities = append(entities, repository.ToRepositoryEntity())
	}

	ctx.JSON(http.StatusOK, entities)
}

func (c Controller) GetLatestUserSyncedRepositories(ctx *gin.Context) {
	repositories, err := c.service.GetLatestUserSyncedRepositories(ctx)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	entities := make([]*domain.RepositoryEntity, 0, len(repositories))
	for _, repository := range repositories {
		entities = append(entities, repository.ToRepositoryEntity())
	}

	ctx.JSON(http.StatusOK, entities)
}
