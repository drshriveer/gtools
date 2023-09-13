package log_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/drshriveer/gtools/log"
)

func TestInitLogger(t *testing.T) {
	f1 := zap.String("f1", "f1")
	f2 := zap.Uint64("f2", 123)
	f3 := zap.Float64("f3", 3.14159)
	f4 := zap.String("f4", "f4")

	core, ob := observer.New(zapcore.WarnLevel)
	zap.ReplaceGlobals(zap.New(core))
	ctx1 := log.InitLogger(context.TODO(), f1)

	// ensure that the initial logger has the common field...
	logErrorAndAssert(ctx1, t, ob, "msg1", f1)

	// InitLogger should overwrite fields.
	ctx2 := log.InitLogger(ctx1, f2)
	logErrorAndAssert(ctx1, t, ob, "msg2", f1)
	logErrorAndAssert(ctx2, t, ob, "msg3", f2)

	// WithField should propagate the field.
	ctx3 := log.WithFields(ctx2, f3)
	logErrorAndAssert(ctx1, t, ob, "msg4", f1)
	logErrorAndAssert(ctx2, t, ob, "msg5", f2, f3)
	logErrorAndAssert(ctx3, t, ob, "msg6", f2, f3)

	// Child logger should not propagate.
	ctx4 := log.ChildLogger(ctx3, f4)
	logErrorAndAssert(ctx1, t, ob, "msg7", f1)
	logErrorAndAssert(ctx2, t, ob, "msg8", f2, f3)
	logErrorAndAssert(ctx3, t, ob, "msg9", f2, f3)
	logErrorAndAssert(ctx4, t, ob, "msg10", f2, f3, f4)

	// Ensure we don't emit infos...
	log.Log(ctx4).Info("msg11")
	assert.Empty(t, ob.TakeAll())

	// Set Level to info and ensure we get the log.
	log.SetLevel(ctx4, zapcore.InfoLevel)
	logLevelAndAssert(ctx4, t, zapcore.InfoLevel, ob, "msg12", f2, f3, f4)

	// Ensure we still don't emit debug.
	log.Log(ctx4).Debug("msg13")
	assert.Empty(t, ob.TakeAll())

	// Last ensure setting debug works:
	log.EnableDebug(ctx4)
	logLevelAndAssert(ctx4, t, zapcore.DebugLevel, ob, "msg14", f2, f3, f4)
}

func logErrorAndAssert(
	ctx context.Context, t *testing.T, logs *observer.ObservedLogs, msg string, fields ...zap.Field,
) {
	log.Log(ctx).Error(msg)
	assertLog(t, logs, msg, fields...)
}

func logLevelAndAssert(
	ctx context.Context,
	t *testing.T,
	level zapcore.Level,
	logs *observer.ObservedLogs,
	msg string,
	fields ...zap.Field,
) {
	log.Log(ctx).Log(level, msg)
	assertLog(t, logs, msg, fields...)
}

func assertLog(t *testing.T, logs *observer.ObservedLogs, msg string, fields ...zap.Field) {
	lls := logs.TakeAll()
	require.Len(t, lls, 1)
	ll := lls[0]
	assert.Equal(t, msg, ll.Message)
	assert.Len(t, ll.Context, len(fields))
	assert.ElementsMatch(t, ll.Context, fields)
}
