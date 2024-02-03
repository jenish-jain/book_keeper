package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

var log *slog.Logger

type Logger interface {
	Debug(log string, args ...any)
	DebugWithCtx(ctx context.Context, logMsg string, args ...any)
	Info(log string, args ...any)
	Error(log string, args ...any)
	ErrorWithCtx(ctx context.Context, logMsg string, args ...any)
	Warn(log string, args ...any)
}

type Impl struct {
}

/*
Debug log
*/
func Debug(logMsg string, args ...any) {
	log.Debug(logMsg, args...)
}

func DebugWithCtx(ctx context.Context, logMsg string, args ...any) {
	args = append(args, "request_id", ctx.Value("request_id"))
	log.DebugContext(ctx, logMsg, args...)
}

/*
Info log
*/
func Info(logMsg string, args ...any) {
	log.Info(logMsg, args...)
}

/*
Warn log
*/
func Warn(logMsg string, args ...any) {
	log.Warn(logMsg, args...)
}

/*
Error log
*/
func Error(logMsg string, args ...any) {
	log.Error(logMsg, args...)
}

func ErrorWithCtx(ctx context.Context, logMsg string, args ...any) {
	args = append(args, "request_id", ctx.Value("request_id"))
	log.ErrorContext(ctx, logMsg, args...)
}

func Init(logLevel string) {
	handlerOptions := slog.HandlerOptions{Level: getLogLevel(logLevel)}
	log = slog.New(slog.NewJSONHandler(os.Stdout, &handlerOptions))
	slog.SetDefault(log)
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
