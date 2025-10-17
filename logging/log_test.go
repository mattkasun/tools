package logging_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/Kairum-Labs/should"
	"github.com/mattkasun/tools/logging"
)

func TestDefaultLogger(t *testing.T) {
	logger := logging.DefaultLogger()
	should.NotBeNil(t, logger)
	should.BeTrue(t, logger.Handler().Enabled(context.Background(), slog.LevelInfo))
	should.BeFalse(t, logger.Handler().Enabled(context.Background(), slog.LevelDebug))
}

func TestTextLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.TextLogger(logging.Output(&buf))
	should.NotBeNil(t, logger)
	logger.Info("hello world")
	should.ContainSubstring(t, buf.String(), "hello world")
}

func TestJsonLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.JSONLogger(logging.Output(&buf))
	should.NotBeNil(t, logger)

	logger.Info("json test", "foo", "bar")
	var out map[string]any
	err := json.Unmarshal(buf.Bytes(), &out)
	should.BeNil(t, err)
	should.BeEqual(t, "json test", out["msg"])
	should.BeEqual(t, "bar", out["foo"])
}

func TestDicardLogger(t *testing.T) {
	logger := logging.DiscardLogger()
	should.NotBeNil(t, logger)
	should.BeFalse(t, logger.Handler().Enabled(context.Background(), slog.LevelError))
	should.BeFalse(t, logger.Handler().Enabled(context.Background(), slog.LevelWarn))
	should.BeFalse(t, logger.Handler().Enabled(context.Background(), slog.LevelInfo))
	should.BeFalse(t, logger.Handler().Enabled(context.Background(), slog.LevelDebug))
}

func TestTimeFormat(t *testing.T) {
	var buf bytes.Buffer
	format := "2006-01-02"
	logger := logging.TextLogger(
		logging.Output(&buf),
		logging.TimeFormat(format),
	)
	logger.Info("check time")
	should.ContainSubstring(t, buf.String(), time.Now().Format(format))
}

func TestWithSource(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.TextLogger(
		logging.Output(&buf),
		logging.WithSource(),
	)
	logger.Info("with source")
	logOutput := buf.String()
	// Should only contain filename
	should.ContainSubstring(t, logOutput, "source=/")
}

func TestTruncateSource(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.TextLogger(
		logging.Output(&buf),
		logging.TruncateSource(),
	)
	logger.Info("with source")
	logOutput := buf.String()
	// Should only contain filename
	should.ContainSubstring(t, logOutput, "source=log_test.go:")
}

func TestLevelOption(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.JSONLogger(
		logging.Output(&buf),
		logging.Level(slog.LevelWarn),
	)
	logger.Info("should not appear")
	logger.Warn("should appear")
	logOutput := buf.String()
	var out map[string]string
	err := json.Unmarshal(buf.Bytes(), &out)
	should.BeNil(t, err)
	should.NotContainValue(t, out, "should not appear")
	should.ContainSubstring(t, logOutput, "should appear")
}

func TestSetDefaultOption(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		var buf bytes.Buffer
		logger := logging.JSONLogger(
			logging.Output(&buf),
			logging.SetDefault(),
			logging.TimeFormat(time.DateOnly),
		)
		logger.Info("testing", "hello", "world")
		first := buf
		buf.Reset()
		slog.Info("testing", "hello", "world")
		second := buf
		should.BeEqual(t, first, second)
	})
	t.Run("unset", func(t *testing.T) {
		var buf bytes.Buffer
		logger := logging.JSONLogger(
			logging.Output(&buf),
			logging.TimeFormat(time.DateOnly),
		)
		logger.Info("testing", "hello", "world")
		first := buf
		buf.Reset()
		slogger := slog.NewJSONHandler(&buf, nil)
		slog.SetDefault(slog.New(slogger))
		slog.Info("testing", "hello", "world")
		should.NotBeEqual(t, first.String(), buf.String())
	})
}
