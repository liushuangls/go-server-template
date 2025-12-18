package controller

import (
	"github.com/labstack/echo/v4"

	"github.com/liushuangls/go-server-template/dto/request"
	"github.com/liushuangls/go-server-template/routes/common"
)

type HealthRoute struct {
	Options
}

func NewHealthRoute(opt Options) *HealthRoute {
	return &HealthRoute{opt}
}

func (u *HealthRoute) RegisterRoute(router *echo.Group) {
	router.GET("/health", u.health)
}

func (u *HealthRoute) health(c echo.Context) error {
	var req request.HealthReq
	if err := c.Bind(&req); err != nil {
		return err
	}

	return common.WrapResp(c)(u.HealthService.Health(c.Request().Context(), &req))
}
