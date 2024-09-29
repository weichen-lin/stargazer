package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/domain"
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
	crontab, _ := c.db.GetCrontab(ctx)
	if crontab != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "already set crontab"})
		return
	}

	newCrontab := domain.NewCrontab()
	err := c.db.SaveCrontab(ctx, newCrontab)

	if err != nil {
		handleError(err, ctx)
		return
	}

	ctx.JSON(http.StatusCreated, newCrontab.ToCrontabEntity())
}

type UpdateQuery struct {
	Hour int `form:"hour" binding:"required,min=0,max=23"`
}

func getTime(hour int) (time.Time, error) {
	if hour < 0 || hour > 23 {
		return time.Time{}, fmt.Errorf("小時數必須介於 0 到 23 之間")
	}

	now := time.Now()

	year, month, day := now.Date()

	t := time.Date(year, month, day, hour, 0, 0, 0, time.Local)

	return t, nil
}

func (c *Controller) UpdateCrontab(ctx *gin.Context) {
	user, err := c.db.GetUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var query UpdateQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	crontab, err := c.db.GetCrontab(ctx)
	if err != nil {
		handleError(err, ctx)
		return
	}

	parsedTime, err := getTime(query.Hour)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}

	now := time.Now()
	crontab.SetTriggeredAt(parsedTime.Format(time.RFC3339))
	crontab.SetUpdatedAt(now.Format(time.RFC3339))
	crontab.UpdateVersion()

	err = c.db.SaveCrontab(ctx, crontab)

	c.scheduler.Update(user.Email(), parsedTime.Hour())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, crontab.ToCrontabEntity())
}
