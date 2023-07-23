package log

type customLevelCoreWrapper struct {
	zap.Core
	minLevel zap.Level
}

// Enable overwrites the LevelEnabler interface within zap.Core.
func (c *customLevelCoreWrapper) Enabled(l zap.Level) bool {
	return c.minLevel <= l
}
