package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func Test_GetCrontab(t *testing.T){
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/crontab", testController.GetCronTab)

	 gin.SetMode(gin.TestMode)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/crontab", nil)

    r.ServeHTTP(w, req)

	require.Equal(t, w.Code, http.StatusUnauthorized)
}