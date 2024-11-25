package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
)

const (
	reset       = "\033[0m"
	lightGray   = 37
	red         = 31
	yellow      = 33
	green       = 32
	cyan        = 36
	lightBlue   = 94
	lightRed    = 91
	lightYellow = 93
	white       = 97
)

var Logger *slog.Logger

// PrettyHandler is a custom slog.Handler for colorful logs
type PrettyHandler struct {
	h      slog.Handler
	b      *bytes.Buffer
	m      *sync.Mutex
	level  slog.Level // Minimum log level to show
}

// InitLogger initializes the global logger with the PrettyHandler
func InitLogger(level slog.Level) {
	handler := NewPrettyHandler(level)
	Logger = slog.New(handler)
}

// NewPrettyHandler creates a new instance of PrettyHandler with a specific log level
func NewPrettyHandler(level slog.Level) *PrettyHandler {
	b := &bytes.Buffer{}
	return &PrettyHandler{
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Ignore default keys (time, level, message)
				if a.Key == slog.TimeKey || a.Key == slog.LevelKey || a.Key == slog.MessageKey {
					return slog.Attr{}
				}
				return a
			},
		}),
		b:     b,
		m:     &sync.Mutex{},
		level: level,
	}
}

// Enabled indicates whether the handler should log messages at the given level
func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level // Only log messages at or above the configured level
}

// WithAttrs creates a new handler with the given attributes
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m, level: h.level}
}

// WithGroup creates a new handler with the given group name
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{h: h.h.WithGroup(name), b: h.b, m: h.m, level: h.level}
}

// Handle formats and outputs the log record
func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil // Skip logs below the configured level
	}

	// Format time, level, and message
	timeStr := colorize(lightBlue, r.Time.Format("[15:04:05.000]"))
	level := r.Level.String() + ":"
	switch r.Level {
	case slog.LevelDebug:
		level = colorize(lightGray, level)
	case slog.LevelInfo:
		level = colorize(cyan, level)
	case slog.LevelWarn:
		level = colorize(lightYellow, level)
	case slog.LevelError:
		level = colorize(lightRed, level)
	}
	message := colorize(white, r.Message)

	// Capture additional attributes
	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}
	attrJSON, _ := json.Marshal(attrs)

	// Print the formatted log line
	fmt.Printf("%s %s %s %s\n", timeStr, level, message, colorize(lightGray, string(attrJSON)))
	return nil
}

// computeAttrs extracts additional attributes for the log record
func (h *PrettyHandler) computeAttrs(ctx context.Context, r slog.Record) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()
	if err := h.h.Handle(ctx, r); err != nil {
		return nil, err
	}
	var attrs map[string]any
	if err := json.Unmarshal(h.b.Bytes(), &attrs); err != nil {
		return nil, err
	}
	return attrs, nil
}

// colorize wraps a string with the given color code
func colorize(colorCode int, text string) string {
	return fmt.Sprintf("\033[%dm%s%s", colorCode, text, reset)
}
