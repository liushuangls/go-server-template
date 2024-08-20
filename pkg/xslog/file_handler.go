package xslog

import (
	"context"
	"io"
	"log/slog"
	"path"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	levels = []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError}
)

type fileHandler struct {
	level   slog.Level
	loggers map[slog.Level]*slog.JSONHandler
}

func newFileHandler(c *Config) fileHandler {
	f := fileHandler{
		level:   c.Level,
		loggers: map[slog.Level]*slog.JSONHandler{},
	}
	for _, l := range levels {
		var writers []io.Writer
		writers = append(writers, &lumberjack.Logger{
			Filename:   path.Join(c.FileDir, l.String()+".log"),
			MaxSize:    c.MaxSize,
			MaxBackups: c.MaxBackups,
			MaxAge:     c.MaxAge,
			Compress:   c.Compress,
		})
		for _, e := range c.ExtraWriters {
			if l >= e.Level {
				writers = append(writers, e.Writer)
			}
		}
		f.loggers[l] = slog.NewJSONHandler(io.MultiWriter(writers...), &slog.HandlerOptions{
			AddSource:   c.AddSource,
			Level:       c.Level,
			ReplaceAttr: c.ReplaceAttr,
		})
	}
	return f
}

func (f fileHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if level >= f.level {
		return true
	}
	return false
}

func (f fileHandler) Handle(ctx context.Context, record slog.Record) error {
	logger := f.loggers[record.Level]
	if logger == nil {
		return nil
	}
	return logger.Handle(ctx, record)
}

func (f fileHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	logger := f.loggers[f.level]
	if logger == nil {
		return nil
	}
	return logger.WithAttrs(attrs)
}

func (f fileHandler) WithGroup(name string) slog.Handler {
	logger := f.loggers[f.level]
	if logger == nil {
		return nil
	}
	return logger.WithGroup(name)
}
