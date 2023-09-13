package log

import (
	"go.uber.org/zap/zapcore"
)

type customLevelCoreWrapper struct {
	zapcore.Core
	minLevel zapcore.Level
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
