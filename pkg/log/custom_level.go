package log

import (
	"go.uber.org/zap/zapcore"
)

type customLevelCoreWrapper struct {
	zapcore.Core
	minLevel zapcore.Level
}

// Enable overwrites the LevelEnabler interface within zap.Core.
func (c *customLevelCoreWrapper) Enabled(l zapcore.Level) bool {
	return c.minLevel <= l
}
