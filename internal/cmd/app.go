package cmd

import (
	configs "github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/internal/routes"
	"github.com/liushuangls/go-server-template/pkg/jwt"
	"github.com/liushuangls/go-server-template/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	http *routes.HttpEngine
}

func NewApp(http *routes.HttpEngine) *App {
	return &App{http: http}
}

func NewLogger(conf *configs.Config) *zap.SugaredLogger {
	return logger.NewLogger(&conf.Log)
}

func NewJwt(conf *configs.Config) (*jwt.JWT, error) {
	return jwt.NewJWT(&jwt.Config{Issuer: conf.Jwt.Issuer, Secret: conf.Jwt.Secret})
}

func (a *App) Run() error {
	return a.http.Run()
}
