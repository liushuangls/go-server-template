//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	configs "github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/internal/cmd"
	"github.com/liushuangls/go-server-template/internal/data"
	"github.com/liushuangls/go-server-template/internal/routes"
	"github.com/liushuangls/go-server-template/internal/service"
)

func app() (*cmd.App, func(), error) {
	panic(wire.Build(
		configs.InitConfig,
		cmd.NewLogger,
		cmd.NewJwt,
		cmd.NewApp,
		routes.ProviderSet,
		data.ProviderSet,
		service.ProviderSet,
	))
}
