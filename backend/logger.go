package main

import (
	"github.com/weichen-lin/kabaka"
	"go.uber.org/zap"
)

type StarGazerLogger struct {
	logger *zap.Logger
}

func (l *StarGazerLogger) Debug(args *kabaka.LogMessage) {
	l.logger.Debug(args.Message, logMessageToZapFields(args)...)
}

func (l *StarGazerLogger) Info(args *kabaka.LogMessage) {
	l.logger.Info(args.Message, logMessageToZapFields(args)...)
}

func (l *StarGazerLogger) Warn(args *kabaka.LogMessage) {
	l.logger.Warn(args.Message, logMessageToZapFields(args)...)
}

func (l *StarGazerLogger) Error(args *kabaka.LogMessage) {
	l.logger.Error(args.Message, logMessageToZapFields(args)...)
}

func logMessageToZapFields(log *kabaka.LogMessage) []zap.Field {
	return []zap.Field{
		zap.String("topic_name", log.TopicName),
		zap.String("action", string(log.Action)),
		zap.String("message_id", log.MessageID.String()),
		zap.String("message", log.Message),
		zap.String("message_status", string(log.MessageStatus)),
		zap.String("subscriber", log.SubScriber.String()),
		zap.Int64("spend_time", log.SpendTime),
		zap.Time("msg_created_at", log.CreatedAt),
	}
}

func NewStarGazerLogger(logger *zap.Logger) *StarGazerLogger {
	return &StarGazerLogger{
		logger: logger,
	}
}
