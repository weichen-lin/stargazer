package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jOpeartion "github.com/weichen-lin/kafka-service/neo4j"
)


type CronTabSetting struct {
	Hour int `json:"hour" binding:"required" form:"hour"`
}

func UpdateCronTabSetting (c *gin.Context) {
	
	userName, ok := c.Value("userName").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	neo4jDriver, exists := c.Get("neo4jDriver")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neo4j driver not found"})
		return
	}

	_, ok = neo4jDriver.(neo4j.DriverWithContext)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Neo4j driver"})
		return
	}

	var cronTabSetting CronTabSetting

	if err := c.ShouldBindQuery(&cronTabSetting) ; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := neo4jOpeartion.AdjustCrontab(neo4jDriver.(neo4j.DriverWithContext), cronTabSetting.Hour, userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update crontab"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Crontab updated successfully"})
}