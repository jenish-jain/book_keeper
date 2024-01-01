package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Logger interface {
	Debug(log string)
	Info(log string)
	Error(log string)
	Warn(log string)
}

type Impl struct {
	logger *slog.Logger
}

/*
Debug log
*/
func (i *Impl) Debug(log string) {
	i.logger.Debug(log)
}

/*
Info log
*/
func (i *Impl) Info(log string) {
	i.logger.Info(log)
}

/*
Warn log
*/
func (i *Impl) Warn(log string) {
	i.logger.Warn(log)
}

/*
Error log
*/
func (i *Impl) Error(log string) {
	i.logger.Error(log)
}

func Init(logLevel string) *Impl {
	handlerOptions := slog.HandlerOptions{Level: getLogLevel(logLevel)}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &handlerOptions))
	slog.SetDefault(logger)
	loggerInstance := &Impl{
		logger: logger,
	}
	return loggerInstance
}

func getLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
