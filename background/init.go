package background

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/kabaka"
	"github.com/weichen-lin/stargazer/infrastructure"
	"github.com/weichen-lin/stargazer/service"
)

type TopicName string

const (
	SyncUserRepositories TopicName = "sync-user-repositories"
)

type Background struct {
	kabaka  *kabaka.Kabaka
	service *service.Service
	config  *infrastructure.Config
}

func NewBackground(service *service.Service, config *infrastructure.Config) *Background {
	kabaka := kabaka.NewKabaka(nil)

	bg := &Background{
		kabaka:  kabaka,
		service: service,
		config:  config,
	}

	err := kabaka.CreateTopic(string(SyncUserRepositories), bg.SyncUserRepositories)
	if err != nil {
		panic(err)
	}

	return bg
}

func bindJSON(ctx *gin.Context, req interface{}) bool {
	logger := infrastructure.GetLogger()

	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "請求參數錯誤"})
		return false
	}
	return true
}

func bindQuery(ctx *gin.Context, req interface{}) bool {
	logger := infrastructure.GetLogger()

	if err := ctx.ShouldBindQuery(req); err != nil {
		logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "請求參數錯誤"})
		return false
	}
	return true
}
