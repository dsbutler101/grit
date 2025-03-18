// Inspired by https://github.com/dusted-go/logging,
// licensed under [Apache 2.0 license](http://www.apache.org/licenses/LICENSE-2.0)
// copyrighted by Dustin Moris Gorski

package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

const (
	colorReset       = "\033[0m"
	colorRed         = "\033[31m"
	colorGreen       = "\033[32m"
	colorYellow      = "\033[33m"
	colorPurple      = "\033[35m"
	colorCyan        = "\033[36m"
	colorGray        = "\033[90m"
	colorLightRed    = "\033[91m"
	colorLightYellow = "\033[93m"
	colorWhite       = "\033[97m"
)

const (
	ErrorKey = "error"
)

func colorize(color string, str string) string {
	return color + str + colorReset
}

type errorEntry struct {
	Error string `json:"error"`
}

type handler struct {
	out io.Writer

	h slog.Handler
	b *bytes.Buffer
	m *sync.Mutex
}

func newHandler(out io.Writer, opts *slog.HandlerOptions) *handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	b := new(bytes.Buffer)

	return &handler{
		out: out,
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			AddSource:   opts.AddSource,
			Level:       opts.Level,
			ReplaceAttr: dropDefaults(opts.ReplaceAttr),
		}),
		b: b,
		m: new(sync.Mutex),
	}
}

func dropDefaults(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	toRemove := map[string]bool{
		slog.TimeKey:    true,
		slog.LevelKey:   true,
		slog.MessageKey: true,
	}

	return func(groups []string, a slog.Attr) slog.Attr {
		if _, ok := toRemove[a.Key]; ok {
			return slog.Attr{}
		}

		if next != nil {
			return next(groups, a)
		}

		return a
	}
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *handler) Handle(ctx context.Context, record slog.Record) error {
	level := record.Level.String()

	switch record.Level {
	case slog.LevelDebug:
		level = colorize(colorCyan, level)
	case slog.LevelInfo:
		level = colorize(colorGreen, level)
	case slog.LevelWarn:
		level = colorize(colorYellow, level)
	case slog.LevelError:
		level = colorize(colorRed, level)
	}

	h.m.Lock()
	defer func() {
		h.m.Unlock()
		h.b.Reset()
	}()

	err := h.h.Handle(ctx, record)
	if err != nil {
		return fmt.Errorf("calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	err = json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return fmt.Errorf("unmarshaling inner handler's Handle: %w", err)
	}

	errorSection := ""
	if errorField, ok := attrs[ErrorKey]; ok {
		if errorFieldStr, ok := errorField.(string); ok {
			delete(attrs, ErrorKey)
			e := errorEntry{
				Error: errorFieldStr,
			}
			b, err := json.MarshalIndent(e, "", "  ")
			if err != nil {
				return fmt.Errorf("marshaling errorEntry: %w", err)
			}

			errorSection = string(b)
		}
	}

	details := ""
	if len(attrs) > 0 {
		b, err := json.MarshalIndent(attrs, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling attrs: %w", err)
		}

		details = string(b)
	}

	_, err = fmt.Fprintln(
		h.out,
		colorize(colorPurple, record.Time.Format("[15:04:05.999 -07:00]")),
		level,
		colorize(colorWhite, record.Message),
		colorize(colorLightYellow, errorSection),
		colorize(colorGray, details),
	)

	return err
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{
		out: h.out,
		h:   h.h.WithAttrs(attrs),
		b:   h.b,
		m:   h.m,
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		out: h.out,
		h:   h.h.WithGroup(name),
		b:   h.b,
		m:   h.m,
	}
}
