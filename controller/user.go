package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c Controller) GetCrontab(ctx *gin.Context) {
	crontab, err := c.service.GetUserCrontab(ctx)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, crontab.ToCrontabDto())
}
