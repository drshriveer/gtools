package log

import (
	"context"
	"sync/atomic"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

// XXX: consider moving this to a different package.
type logHolderKeyType struct{}

var logHolderKey = logHolderKeyType{}

type logHolder struct {
	atomic.Pointer[zap.Logger]
}

func InitLogger(ctx context.Context, fields ...zap.Field) context.Context {
	lh, ok := getOrDefault(ctx)
	if !ok {
		ctx = context.WithValue(ctx, logHolderKey, lh)
	}
	newLogger := lh.Load().With(fields...)
	lh.Store(newLogger)
	return ctx
}

func ChildLogger(ctx context.Context, fields ...zap.Field) context.Context {
	lh, _ := getOrDefault(ctx)
	newLogger := lh.Load().With(fields...)
	lh = &logHolder{}
	lh.Store(newLogger)
	ctx = context.WithValue(ctx, logHolderKey, lh)
	return ctx
}

func Log(ctx context.Context) *zap.Logger {
	lh, _ := getOrDefault(ctx)
	return lh.Load()
}

func EnableDebug(ctx context.Context) context.Context {
	return SetLevel(ctx, zapcore.DebugLevel)
}

func SetLevel(ctx context.Context, level zapcore.Level) context.Context {
	lh, ok := getOrDefault(ctx)
	logger := lh.Load()
	coreWrapper := &customLevelCoreWrapper{
		Core:     logger.Core(),
		minLevel: level,
	}
	lh.Store(zap.New(coreWrapper))
	if !ok {
		ctx = context.WithValue(ctx, logHolderKey, logger)
	}
	return ctx
}

func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	lh, ok := getOrDefault(ctx)
	logger := lh.Load()
	lh.Store(logger.With(fields...))
	if !ok {
		ctx = context.WithValue(ctx, logHolderKey, logger)
	}
	return ctx
}

func getOrDefault(ctx context.Context) (*logHolder, bool) {
	lh, ok := ctx.Value(logHolderKey).(*logHolder)
	if ok {
		return lh, true
	}
	lh = &logHolder{}
	lh.Store(zap.L())
	return lh, false
}
