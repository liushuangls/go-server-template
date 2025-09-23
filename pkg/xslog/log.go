package xslog

import (
	"io"
	"log/slog"
)

type Config struct {
	FileDir    string `koanf:"FileDir"`
	MaxSize    int    `koanf:"MaxSize"`
	MaxBackups int    `koanf:"MaxBackups"`
	MaxAge     int    `koanf:"MaxAge"`
	Compress   bool   `koanf:"Compress"`
	LocalTime  bool   `koanf:"LocalTime"`
	AddSource  bool   `koanf:"AddSource"`

	Level        slog.Level
	ExtraWriters []ExtraWriter
	ReplaceAttr  func(groups []string, a slog.Attr) slog.Attr
}

type ExtraWriter struct {
	Writer io.Writer
	Level  slog.Level
}

func (c *Config) SetLevel(level slog.Level) *Config {
	c.Level = level
	return c
}

func NewFileSlog(c *Config) *slog.Logger {
	return slog.New(newFileHandler(c))
}
