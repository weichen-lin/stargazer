package main

import (
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/kabaka"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func generateFakeLogMessage(action kabaka.Action, msgStatus kabaka.MessageStatus) *kabaka.LogMessage {
	return &kabaka.LogMessage{
		TopicName:     faker.Name(),
		Action:        action,
		MessageID:     uuid.New(),
		Message:       faker.Sentence(),
		MessageStatus: msgStatus,
		SubScriber:    uuid.New(),
		SpendTime:     12345,
		CreatedAt:     time.Now(),
	}
}

func TestStarGazerLogger(t *testing.T) {

	core, recorded := observer.New(zapcore.DebugLevel)

	testLogger := zap.New(core)

	starLogger := NewStarGazerLogger(testLogger)

	testCases := []struct {
		level   string
		message string
		action  kabaka.Action
		status  kabaka.MessageStatus
	}{
		{"debug", "Debug message", "subscribe", "success"},
		{"info", "Info message", "subscribe", "retry"},
		{"warn", "Warn message", "publish", "error"},
		{"error", "Error message", "consume", "success"},
		{"debug", "Debug message", "publish", "retry"},
		{"info", "Info message", "subscribe", "error"},
		{"warn", "Warn message", "subscribe", "success"},
		{"error", "Error message", "subscribe", "retry"},
		{"debug", "Debug message", "consume", "error"},
		{"info", "Info message", "consume", "retry"},
		{"warn", "Warn message", "consume", "error"},
		{"error", "Error message", "publish", "error"},
	}

	for _, tc := range testCases {
		t.Run(tc.level, func(t *testing.T) {
			logMsg := generateFakeLogMessage(tc.action, tc.status)

			switch tc.level {
			case "debug":
				starLogger.Debug(logMsg)
			case "info":
				starLogger.Info(logMsg)
			case "warn":
				starLogger.Warn(logMsg)
			case "error":
				starLogger.Error(logMsg)
			}

			logs := recorded.All()

			require.Equal(t, 1, len(logs))

			log := logs[0]
			contextMap := log.ContextMap()

			require.Equal(t, tc.level, log.Level.String())
			require.Equal(t, logMsg.Message, contextMap["message"])

			require.Equal(t, logMsg.TopicName, contextMap["topic_name"])
			require.Equal(t, string(logMsg.Action), contextMap["action"])
			require.Equal(t, logMsg.MessageID.String(), contextMap["message_id"])
			require.Equal(t, string(logMsg.MessageStatus), contextMap["message_status"])
			require.Equal(t, logMsg.SubScriber.String(), contextMap["subscriber"])
			require.Equal(t, logMsg.SpendTime, contextMap["spend_time"])

			loggedTime, ok := contextMap["msg_created_at"].(time.Time)
			require.True(t, ok, "created_at should be a time.Time")
			require.True(t, loggedTime.Equal(logMsg.CreatedAt), "Logged time should match the original time")

			recorded.TakeAll()
		})
	}
}
