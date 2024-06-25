package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CronTabSetting struct {
	Hour int `json:"hour" binding:"required" form:"hour"`
}

func (c *Controller) UpdateCronTabSetting(ctx *gin.Context) {
	email, ok := ctx.Value("email").(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var cronTabSetting CronTabSetting

	if err := ctx.ShouldBindQuery(&cronTabSetting); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.db.UpdateCrontab(cronTabSetting.Hour, email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update crontab"})
		return
	}

	err = c.scheduler.Update(email, cronTabSetting.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update crontab"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Crontab updated successfully"})
}

// func (c *Controller) GetAllCronJobs(ctx *gin.Context) {
// 	jobs := c.scheduler.GetAll()
// 	ctx.JSON(http.StatusOK, gin.H{"jobs": jobs})
// }
