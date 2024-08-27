package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/kafka-service/db"
	"github.com/weichen-lin/kafka-service/domain"
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

func (c *Controller) GetCrontab(ctx *gin.Context) {
	crontab, err := c.db.GetCrontab(ctx)
	if err != nil {
		handleError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, crontab.ToCrontabEntity())
}

func (c *Controller) CreateCrontab(ctx *gin.Context) {
	crontab, err := c.db.GetCrontab(ctx)
	if crontab != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "already set crontab"})
		return
	}

	newCrontab := domain.NewCrontab()
	err = c.db.CreateCrontab(ctx, newCrontab)
	fmt.Println(err)
	if err != nil {
		handleError(err, ctx)
		return
	}

	ctx.JSON(http.StatusCreated, newCrontab.ToCrontabEntity())
}