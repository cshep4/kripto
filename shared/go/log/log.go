package log

import (
	"context"
	"os"

	"github.com/palantir/witchcraft-go-logging/wlog"
	wlogzap "github.com/palantir/witchcraft-go-logging/wlog-zap"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

type Param = svc1log.Param

func New(level string) svc1log.Logger {
	wlog.SetDefaultLoggerProvider(wlogzap.LoggerProvider())
	return svc1log.New(os.Stdout, wlog.LogLevel(level))
}

func WithServiceName(ctx context.Context, logger svc1log.Logger, service string) context.Context {
	return svc1log.WithLoggerParams(
		svc1log.WithLogger(ctx, logger),
		svc1log.SafeParam("service", service),
	)
}

func WithFunctionName(ctx context.Context, logger svc1log.Logger, service, function string) context.Context {
	return svc1log.WithLoggerParams(
		svc1log.WithLogger(ctx, logger),
		svc1log.SafeParam("service", service),
		svc1log.SafeParam("function", function),
	)
}

func SafeParam(k string, v interface{}) Param {
	return svc1log.SafeParam(k, v)
}

func ErrorParam(err error) Param {
	return svc1log.Stacktrace(err)
}

func Debug(ctx context.Context, msg string, params ...Param) {
	svc1log.FromContext(ctx).Debug(msg, params...)
}

func Info(ctx context.Context, msg string, params ...Param) {
	svc1log.FromContext(ctx).Info(msg, params...)
}

func Warn(ctx context.Context, msg string, params ...Param) {
	svc1log.FromContext(ctx).Warn(msg, params...)
}

func Error(ctx context.Context, msg string, params ...Param) {
	svc1log.FromContext(ctx).Error(msg, params...)
}

func Fatal(ctx context.Context, msg string, params ...Param) {
	Error(ctx, msg, params...)
	os.Exit(1)
}
