package log

import (
	"context"
	"sync/atomic"
	"testing"

	"go.uber.org/zap/zaptest"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

// XXX: consider moving this to a different package.
type logHolderKeyType struct{}

var logHolderKey = logHolderKeyType{}

type logHolder struct {
	atomic.Pointer[zap.Logger]
}

// InitLogger with fields supplied (if so supplied). If a context already has a logger,
// the fields will be added to it.
func InitLogger(ctx context.Context, fields ...zap.Field) context.Context {
	lh, ok := getOrDefault(ctx)
	if !ok {
		ctx = context.WithValue(ctx, logHolderKey, lh)
	}
	newLogger := lh.Load().With(fields...)
	lh.Store(newLogger)
	return ctx
}

// ChildLogger _always_ creates a new logger with the fields of the parent.
// This will prevent propagation of fields higher up the stack so is useful when
// spawning child routines which may write the same field name but be meaningfully different.
func ChildLogger(ctx context.Context, fields ...zap.Field) context.Context {
	lh, _ := getOrDefault(ctx)
	newLogger := lh.Load().With(fields...)
	lh = &logHolder{}
	lh.Store(newLogger)
	ctx = context.WithValue(ctx, logHolderKey, lh)
	return ctx
}

// Log returns the underlying logger with is latest fields.
func Log(ctx context.Context) *zap.Logger {
	lh, _ := getOrDefault(ctx)
	return lh.Load()
}

// TestContext returns a context with a zaptest.Logger tied to a test object.
func TestContext(t *testing.T) context.Context {
	lh := &logHolder{}
	lh.Store(zaptest.NewLogger(t))
	return context.WithValue(context.TODO(), logHolderKey, lh)
}

// EnableDebug turns on debug logs for this context.
func EnableDebug(ctx context.Context) context.Context {
	return SetLevel(ctx, zapcore.DebugLevel)
}

// SetLevel sets a specific log level for this context.
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

// WithFields adds log fields to a context. These will propagate to the nearest initialized or
// child logger.
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
