package logger

import (
	"log/slog"
	"os"
)

func New(opts ...Option) *slog.Logger {
	o := setupOptions(opts)

	hOpts := &slog.HandlerOptions{
		Level:     o.leveler,
		AddSource: o.addSource,
	}

	var h slog.Handler
	if o.customLogFormat {
		h = newHandler(os.Stderr, hOpts)
	} else {
		h = slog.NewJSONHandler(os.Stderr, hOpts)
	}

	return slog.New(h)
}
