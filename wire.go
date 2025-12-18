//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	cmd2 "github.com/liushuangls/go-server-template/app"
	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/crontab"
	"github.com/liushuangls/go-server-template/data"
	"github.com/liushuangls/go-server-template/routes"
	"github.com/liushuangls/go-server-template/service"
)

func app(ctx context.Context) (*cmd2.App, func(), error) {
	panic(wire.Build(
		configs.InitConfig,
		routes.ProviderSet,
		data.ProviderSet,
		service.ProviderSet,
		crontab.ProviderSet,
		cmd2.ProviderSet,
	))
}
