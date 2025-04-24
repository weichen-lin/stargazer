package infrastructure

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type ServiceError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`

	InternalMessage string `json:"internal_message,omitempty"`
}

func (e *ServiceError) Error() string {
	return e.Message
}

func NewServiceError(ctx context.Context, code int, message, internalMessage string) *ServiceError {
	traceId, ok := ctx.Value("traceId").(string)
	if !ok {
		traceId = "unknown traceId"
	}

	logger := GetLogger()
	logger.Error(internalMessage, zap.String("trace-id", traceId), zap.String("message", message))

	fmt.Printf("traceId: %s, message: %s, internalMessage %s", traceId, message, internalMessage)

	err := &ServiceError{
		Code:            code,
		Message:         message,
		InternalMessage: internalMessage,
	}

	return err
}
