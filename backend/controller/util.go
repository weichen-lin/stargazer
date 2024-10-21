package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OtelHeaders struct {
	ginHeader http.Header
}

func NewHeaders(c *gin.Context) *OtelHeaders {
    return &OtelHeaders{
        ginHeader: c.Request.Header,
    }
}

func (o *OtelHeaders) Get(key string) string {
    return o.ginHeader.Get(key)
}

func (o *OtelHeaders) Set(key, value string) {
    o.ginHeader.Set(key, value)
}

func (o *OtelHeaders) Keys() []string {
    keys := make([]string, 0, len(o.ginHeader))
    for key := range o.ginHeader {
        keys = append(keys, key)
    }
    return keys
}