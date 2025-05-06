package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/stargazer/infrastructure"
	"github.com/weichen-lin/stargazer/service"
	"go.uber.org/zap"
)

type Controller struct {
	service *service.Service
}

func NewController(service *service.Service) *Controller {
	return &Controller{
		service: service,
	}
}

func parseServiceError(err error) (int, string) {
	serviceErr, ok := err.(*infrastructure.ServiceError)
	if !ok {
		logger := infrastructure.GetLogger()
		logger.Error("not service error", zap.String("message", err.Error()))
		return http.StatusInternalServerError, "發生非預期錯誤，請稍後再試"
	}

	return serviceErr.Code, serviceErr.Message
}

func handleServiceError(ctx *gin.Context, err error) {
	code, message := parseServiceError(err)
	ctx.JSON(code, gin.H{"error": message})
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
