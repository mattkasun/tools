// Package logging provides slog helpers.
//
//nolint:exhaustruct
package logging

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const (
	kindDefault = "default"
	kindText    = "text"
	kindJSON    = "json"
)

// Logger represents a new slog logger.
type Logger struct {
	*slog.Logger

	kind           string // []string{"default", "text", "json"}
	level          slog.Leveler
	includeSource  bool
	truncateSource bool
	timeFormat     string // loyout constant from time package
	output         io.Writer
}

// Option function.
type Option func(*Logger)

// Level sets the level option.
func Level(level slog.Leveler) Option {
	return func(l *Logger) {
		l.level = level
	}
}

// WithSource sets the includeSource option.
func WithSource() Option {
	return func(l *Logger) {
		l.includeSource = true
		log.SetFlags(log.LstdFlags | log.Llongfile)
	}
}

// TruncateSource sets the includeSource and trucateSource optionsa.
func TruncateSource() Option {
	return func(l *Logger) {
		l.includeSource = true
		l.truncateSource = true
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
}

// TimeFormat set the timeFormat Option.
// has no impact on DefaultLogger.
func TimeFormat(t string) Option {
	return func(l *Logger) {
		l.timeFormat = t
	}
}

// Output sets the output option.
func Output(w io.Writer) Option {
	return func(l *Logger) {
		l.output = w
		log.SetOutput(w)
	}
}

// new returns a new slog logger.
func (l *Logger) new(opts ...Option) *Logger {
	for _, opt := range opts {
		opt(l)
	}
	l.defaults()
	options := &slog.HandlerOptions{Level: l.level}
	if l.includeSource {
		options.AddSource = true
	}
	options.ReplaceAttr = l.replace()
	switch l.kind {
	case kindText:
		l.Logger = slog.New(slog.NewTextHandler(l.output, options))
	case kindJSON:
		l.Logger = slog.New(slog.NewJSONHandler(l.output, options))
	default:
		l.Logger = slog.Default()
	}
	return l
}

func (l *Logger) replace() func(group []string, a slog.Attr) slog.Attr {
	var repSource func(groups []string, a slog.Attr) slog.Attr
	if l.truncateSource {
		repSource = func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.SourceKey {
				if source, _ := attr.Value.Any().(*slog.Source); source != nil {
					source.File = filepath.Base(source.File)
				}
			}
			return attr
		}
	}
	repTime := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			t := a.Value.Time()
			a.Value = slog.StringValue(t.Format(l.timeFormat))
		}
		return a
	}

	return func(groups []string, a slog.Attr) slog.Attr {
		a = repTime(groups, a)
		if l.truncateSource {
			a = repSource(groups, a)
		}
		return a
	}
}

// DefaultLogger returns a logger with slog default handler.
func DefaultLogger(opts ...Option) *Logger {
	l := &Logger{kind: kindDefault}
	return l.new(opts...)
}

// TextLogger returns a logger with a slog.TextHandler.
func TextLogger(opts ...Option) *Logger {
	l := &Logger{kind: kindText}
	return l.new(opts...)
}

// JSONLogger returns a logger with a slog.JSONHandler.
func JSONLogger(opts ...Option) *Logger {
	l := Logger{kind: kindJSON}
	return l.new(opts...)
}

// DiscardLogger returns a logger with slog.DiscardHandler.
func DiscardLogger() *Logger {
	return &Logger{
		Logger: slog.New(slog.DiscardHandler),
	}
}

func (l *Logger) defaults() {
	if l.level == nil {
		l.level = slog.LevelInfo
	}
	if l.timeFormat == "" {
		l.timeFormat = time.DateTime
	}
	if l.output == nil {
		l.output = os.Stderr
	}
}
