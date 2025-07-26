// Package logging provides slog helpers.
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

type LogOption func(*Logger)

func Level(level slog.Leveler) LogOption {
	return func(l *Logger) {
		l.level = level
	}
}

func WithSource() LogOption {
	return func(l *Logger) {
		l.includeSource = true
	}
}

func TruncateSource() LogOption {
	return func(l *Logger) {
		l.includeSource = true
		l.truncateSource = true
	}
}

func TimeFormat(t string) LogOption {
	return func(l *Logger) {
		l.timeFormat = t
	}
}

func Output(w io.Writer) LogOption {
	return func(l *Logger) {
		l.output = w
	}
}

// New returns a new slog logger.
func (l *Logger) new(opts ...LogOption) *Logger { //nolint:cyclop
	for _, opt := range opts {
		opt(l)
	}
	l.defaults()
	var flag int
	var repSource func(groups []string, a slog.Attr) slog.Attr
	if l.includeSource {
		flag = log.Llongfile
	}
	if l.truncateSource {
		flag = log.Lshortfile
		repSource = func(_ []string, a slog.Attr) slog.Attr { //nolint:varnamelen
			if a.Key == slog.SourceKey {
				if source, _ := a.Value.Any().(*slog.Source); source != nil {
					source.File = filepath.Base(source.File)
				}
			}
			return a
		}
	}
	repTime := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			t := a.Value.Time()
			a.Value = slog.StringValue(t.Format(l.timeFormat))
		}
		return a
	}
	log.SetFlags(log.LstdFlags | flag)
	options := &slog.HandlerOptions{Level: l.level} //nolint:exhaustruct
	if l.includeSource {
		options.AddSource = true
	}
	replace := func(groups []string, a slog.Attr) slog.Attr {
		a = repTime(groups, a)
		if l.truncateSource {
			a = repSource(groups, a)
		}
		return a
	}
	options.ReplaceAttr = replace
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

func DefaultLogger(opts ...LogOption) *Logger {
	l := &Logger{kind: kindDefault}
	return l.new(opts...)
}

func TextLogger(opts ...LogOption) *Logger {
	l := &Logger{kind: kindText}
	return l.new(opts...)
}

func JsonLogger(opts ...LogOption) *Logger {
	l := Logger{kind: kindJSON}
	return l.new(opts...)
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
