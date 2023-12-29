package crontab

import (
	"github.com/google/wire"

	"github.com/liushuangls/go-server-template/internal/clients/logdbsync"
	"github.com/liushuangls/go-server-template/internal/clients/serverconf"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Options), "*"),
	NewClient,
)

type Options struct {
	LogDBSync  *logdbsync.Client
	ServerConf *serverconf.ServerConf
}
