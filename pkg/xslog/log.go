package xslog

import (
	"io"
	"log/slog"
)

type Config struct {
	FileDir    string `yaml:"FileDir"`
	MaxSize    int    `yaml:"MaxSize"`
	MaxBackups int    `yaml:"MaxBackups"`
	MaxAge     int    `yaml:"MaxAge"`
	Compress   bool   `yaml:"Compress"`
	LocalTime  bool   `yaml:"LocalTime"`
	AddSource  bool   `yaml:"AddSource"`

	Level       slog.Level
	ExtraWriter io.Writer
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
}

func (c *Config) SetLevel(level slog.Level) *Config {
	c.Level = level
	return c
}

func NewFileSlog(c *Config) *slog.Logger {
	return slog.New(newFileHandler(c))
}
