package common

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/liushuangls/go-server-template/pkg/ecode"
)

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    any    `json:"data"`
}

func NewResp(data interface{}, err error) (int, *Resp) {
	var (
		code     = 0
		msg      = ""
		httpCode = http.StatusOK
	)
	ec := ecode.FromError(err)
	if ec != nil {
		code = ec.Code
		msg = ec.Message
		httpCode = ec.HttpCode
	}
	if code == ecode.UnknownCode {
		msg = "Internal Server Error"
		slog.Error("NewResp receive unknown error", "err", ec)
	}
	return httpCode, &Resp{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

func WrapResp(c echo.Context) func(data any, err error) error {
	return func(data any, err error) error {
		if err != nil {
			return err
		}

		return c.JSON(NewResp(data, nil))
	}
}
