package logger

import (
	"log/slog"
)

const (
	defaultAddSource = false
	customLogFormat  = false
)

type Option func(o *options)

type options struct {
	addSource       bool
	customLogFormat bool
	leveler         slog.Leveler
}

func setupOptions(opts []Option) options {
	l := &slog.LevelVar{}
	l.Set(slog.LevelInfo)

	o := options{
		addSource:       defaultAddSource,
		customLogFormat: customLogFormat,
		leveler:         l,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

func WithCustomLogFormat(customLogFormat bool) Option {
	return func(o *options) {
		o.customLogFormat = customLogFormat
	}
}

func WithAddSource(addSource bool) Option {
	return func(o *options) {
		o.addSource = addSource
	}
}

func WithLeveler(level slog.Leveler) Option {
	return func(o *options) {
		o.leveler = level
	}
}
