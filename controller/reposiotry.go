package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c Controller) GetLanguageDistribution(ctx *gin.Context) {
	languageDistribution, err := c.service.GetLanguageDistribution(ctx)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, languageDistribution)
}
