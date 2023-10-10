package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// customLevelCoreWrapper wraps a zap core to enable setting any log level for a specific logger.
type customLevelCoreWrapper struct {
	zapcore.Core

	// minLevel to log.
	minLevel zapcore.Level
}

// CustomLevelLogger will enable any log level on a given logger.
func CustomLevelLogger(logger *zap.Logger, level zapcore.Level) *zap.Logger {
	return logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &customLevelCoreWrapper{
			Core:     core,
			minLevel: level,
		}
	}))
}

// Level returns the level of this wrapped core.
func (c *customLevelCoreWrapper) Level() zapcore.Level {
	return c.minLevel
}

// Enabled overwrites the LevelEnabler interface within zapcore.Core.
func (c *customLevelCoreWrapper) Enabled(l zapcore.Level) bool {
	return c.minLevel <= l
}

// Check determines whether the supplied Entry should be logged (using the
// embedded LevelEnabler and possibly some extra logic). If the entry
// should be logged, the Core adds itself to the CheckedEntry and returns
// the result.
//
// Callers must use Check before calling Write.
func (c *customLevelCoreWrapper) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}
