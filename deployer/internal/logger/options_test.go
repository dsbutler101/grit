package logger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCustomLeveler struct{}

func (t *testCustomLeveler) Level() slog.Level {
	return slog.LevelDebug
}

func TestSetupOptions(t *testing.T) {
	o := setupOptions(nil)

	assert.Equal(t, defaultAddSource, o.addSource)
	assert.Equal(t, defaultCustomLogFormat, o.customLogFormat)
	assert.IsType(t, &slog.LevelVar{}, o.leveler)

	o2 := setupOptions([]Option{
		WithAddSource(true),
		WithCustomLogFormat(true),
		WithLeveler(new(testCustomLeveler)),
	})

	assert.Equal(t, true, o2.addSource)
	assert.Equal(t, true, o2.customLogFormat)
	assert.IsType(t, &testCustomLeveler{}, o2.leveler)
}
