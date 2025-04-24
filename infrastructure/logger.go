package infrastructure

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func InitLogger(level zapcore.Level, development bool) {
	once.Do(func() {
		var err error
		config := zap.Config{
			Level:             zap.NewAtomicLevelAt(level),
			Development:       development,
			Encoding:          "json",
			EncoderConfig:     zap.NewProductionEncoderConfig(),
			OutputPaths:       []string{"stdout"},
			ErrorOutputPaths:  []string{"stderr"},
			DisableStacktrace: true,
		}

		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		logger, err = config.Build()

		if err != nil {
			panic(err)
		}
	})
}

func GetLogger() *zap.Logger {
	if logger == nil {
		panic("logger not initialized. Call InitLogger first.")
	}
	return logger
}

func GetSugarLogger() *zap.SugaredLogger {
	return GetLogger().Sugar()
}
