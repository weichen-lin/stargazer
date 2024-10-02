package controller

import (
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
	TriggeredAt time.Time `form:"triggered_at" json:"time" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
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

	now := time.Now()

	crontab.SetTriggeredAt(query.TriggeredAt.Format(time.RFC3339))
	crontab.SetUpdatedAt(now.Format(time.RFC3339))
	crontab.UpdateVersion()

	err = c.db.SaveCrontab(ctx, crontab)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fn := func() error {
		err := c.kabaka.Publish("star-syncer", []byte(`{"email":"`+user.Email()+`","page":1}`))
		return err
	}

	err = c.scheduler.Update(user.Email(), query.TriggeredAt, fn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, crontab.ToCrontabEntity())
}
