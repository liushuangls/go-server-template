package common

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/liushuangls/go-server-template/pkg/ecode"
)

func EchoErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var he *echo.HTTPError
	if errors.As(err, &he) {
		err = ecode.New(1000, he.Code, fmt.Sprintf("%v", he.Message))
	}

	_ = c.JSON(NewResp(nil, err))
}
