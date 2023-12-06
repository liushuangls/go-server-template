package xslog

import (
	"log/slog"
)

func Error(err error) slog.Attr {
	return slog.String("err", err.Error())
}
