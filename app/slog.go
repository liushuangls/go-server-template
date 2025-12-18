package cmd

import (
	"log/slog"
	"os"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/pkg/xslog"
)

func NewDefaultSlog(conf *configs.Config) *slog.Logger {
	var extraWriters []xslog.ExtraWriter

	if conf.IsDebugMode() {
		conf.Log.Level = slog.LevelDebug
		extraWriters = append(extraWriters, xslog.ExtraWriter{
			Writer: os.Stdout,
			Level:  slog.LevelDebug,
		})
	}

	conf.Log.ExtraWriters = extraWriters
	fileLogger := xslog.NewFileSlog(&conf.Log)
	slog.SetDefault(fileLogger)

	return fileLogger
}
