package zapper

import "go.uber.org/zap"

var tops = []TeeOption{
	{
		Filename: "logstash/logs/access.log",
		Ropt: RotateOptions{
			MaxSize:    1,
			MaxAge:     1,
			MaxBackups: 3,
			Compress:   true,
		},
		Lef: func(lvl Level) bool {
			return lvl <= InfoLevel
		},
	},
	{
		Filename: "logstash/logs/error.log",
		Ropt: RotateOptions{
			MaxSize:    1,
			MaxAge:     1,
			MaxBackups: 3,
			Compress:   true,
		},
		Lef: func(lvl Level) bool {
			return lvl > InfoLevel
		},
	},
}

var std = NewTeeWithRotate(tops)

var (
	Skip       = zap.Skip
	Binary     = zap.Binary
	Bool       = zap.Bool
	Boolp      = zap.Boolp
	ByteString = zap.ByteString
	Float64    = zap.Float64
	Float64p   = zap.Float64p
	Float32    = zap.Float32
	Float32p   = zap.Float32p
	Durationp  = zap.Durationp
	Int        = zap.Int
	Intp       = zap.Intp
	String     = zap.String

	Info  = std.Info
	Warn  = std.Warn
	Error = std.Error
	Debug = std.Debug
)

type Logger struct {
	l     *zap.Logger
}

type Field = zap.Field

func (l *Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

func Sync() error {
	if std != nil {
		return std.Sync()
	}
	return nil
}
