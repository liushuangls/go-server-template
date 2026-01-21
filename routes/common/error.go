package common

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/liushuangls/go-server-template/pkg/ecode"
)

func EchoErrorHandler(c *echo.Context, err error) {
	if r, _ := echo.UnwrapResponse(c.Response()); r != nil && r.Committed {
		return
	}

	var (
		sErr *ecode.Error
		eErr *echo.HTTPError
		sc   echo.HTTPStatusCoder
	)
	if errors.As(err, &sErr) {
	} else if errors.As(err, &eErr) {
		code := eErr.StatusCode()
		if code >= 400 && code < 500 {
			sErr = ecode.InvalidParams.WithHttpCodeCause(code, eErr.Unwrap())
		} else {
			sErr = ecode.New(ecode.UnknownCode, code, http.StatusText(code)).WithCause(eErr.Unwrap())
		}
	} else if errors.As(err, &sc) {
		code := sc.StatusCode()
		if code >= 400 && code < 500 {
			sErr = ecode.InvalidParams.WithHttpCodeCause(code, err)
		} else {
			sErr = ecode.New(ecode.UnknownCode, code, http.StatusText(code)).WithCause(err)
		}
	} else {
		sErr = ecode.InternalServerErr.WithCause(fmt.Errorf("%s received unknown error: %w", c.Path(), err))
	}

	_ = c.JSON(NewResp(nil, sErr))
}
