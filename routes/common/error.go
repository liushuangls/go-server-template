package common

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v5"

	"github.com/liushuangls/go-server-template/pkg/ecode"
)

func EchoErrorHandler(c *echo.Context, err error) {
	echoResp, unErr := echo.UnwrapResponse(c.Response())
	if unErr != nil {
		unErr = ecode.InternalServerErr.WithCause(fmt.Errorf("EchoErrorHandler: UnwrapResponse err: %w", unErr))
		_ = c.JSON(NewResp(nil, unErr))
		return
	}
	if echoResp.Committed {
		return
	}

	var he *echo.HTTPError
	if errors.As(err, &he) {
		eErr := ecode.FromError(he.Unwrap())
		if eErr == nil || eErr.Code == ecode.UnknownCode {
			err = ecode.New(1000, he.Code, fmt.Sprintf("%v", he.Message))
		}
	}

	_ = c.JSON(NewResp(nil, err))
}
