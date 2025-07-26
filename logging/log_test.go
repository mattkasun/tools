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

func TestDefaultLoggerCreation(t *testing.T) {
	t.Helper()
	logger := logging.DefaultLogger()
	should.NotBeNil(t, logger)
	should.BeTrue(t, logger.Handler().Enabled(context.TODO(), slog.LevelInfo))
}

func TestTextLoggerCreation(t *testing.T) {
	t.Helper()
	var buf bytes.Buffer
	logger := logging.TextLogger(logging.Output(&buf))
	should.NotBeNil(t, logger)

	logger.Info("hello world")
	should.ContainSubstring(t, buf.String(), "hello world")
}

func TestJsonLoggerCreation(t *testing.T) {
	t.Helper()
	var buf bytes.Buffer
	logger := logging.JsonLogger(logging.Output(&buf))
	should.NotBeNil(t, logger)

	logger.Info("json test", "foo", "bar")

	var out map[string]any
	err := json.Unmarshal(buf.Bytes(), &out)
	should.BeNil(t, err)
	should.BeEqual(t, "json test", out["msg"])
	should.BeEqual(t, "bar", out["foo"])
}

func TestTimeFormat(t *testing.T) {
	t.Helper()
	var buf bytes.Buffer
	tf := "2006-01-02"
	logger := logging.TextLogger(
		logging.Output(&buf),
		logging.TimeFormat(tf),
	)

	logger.Info("check time")
	should.ContainSubstring(t, buf.String(), time.Now().Format(tf))
}

func TestWithSource(t *testing.T) {
	t.Helper()
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
	t.Helper()
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
	t.Helper()
	var buf bytes.Buffer
	logger := logging.JsonLogger(
		logging.Output(&buf),
		logging.Level(slog.LevelWarn),
	)

	logger.Info("should not appear")
	logger.Warn("should appear")
	logOutput := buf.String()
	var out map[string]string
	err := json.Unmarshal(buf.Bytes(), &out)
	should.BeNil(t, err)
	should.NotContainKey(t, out, "should appear")
	should.ContainSubstring(t, logOutput, "should appear")
}
