package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/kafka-service/db"
)

func handleError(err error, ctx *gin.Context) {
	switch err {
	case db.ErrNotFoundEmailAtContext:
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UnAuthorized"})
	case db.ErrNotFoundCrontab:
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server is not response now"})
	}
}

func (c *Controller) GetCronTab(ctx *gin.Context) {
	crontab, err := c.db.GetCrontab(ctx)
	if err != nil {
		handleError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, crontab.ToCrontabEntity())
}
