package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/kafka-service/db"
)

type CronTabSetting struct {
	Hour int `json:"hour" binding:"required" form:"hour"`
}

func UpdateCronTabSetting(c *gin.Context) {
	email, ok := c.Value("email").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	driver, exists := c.Get("neo4jDriver")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neo4j driver not found"})
		return
	}

	driver, ok = driver.(db.Database)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Neo4j driver"})
		return
	}

	var cronTabSetting CronTabSetting

	if err := c.ShouldBindQuery(&cronTabSetting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := driver.(db.Database).UpdateCrontab(cronTabSetting.Hour, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update crontab"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Crontab updated successfully"})
}
