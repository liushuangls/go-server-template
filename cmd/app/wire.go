//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/internal/cmd"
	"github.com/liushuangls/go-server-template/internal/crontab"
	"github.com/liushuangls/go-server-template/internal/data"
	"github.com/liushuangls/go-server-template/internal/routes"
	"github.com/liushuangls/go-server-template/internal/service"
)

func app(ctx context.Context) (*cmd.App, func(), error) {
	panic(wire.Build(
		configs.InitConfig,
		routes.ProviderSet,
		data.ProviderSet,
		service.ProviderSet,
		crontab.ProviderSet,
		cmd.ProviderSet,
	))
}
