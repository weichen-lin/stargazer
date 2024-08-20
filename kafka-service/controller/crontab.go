package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) GetCronTab(ctx *gin.Context) {
	crontab, err := c.db.GetCrontab(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, crontab.ToCrontabEntity())
}