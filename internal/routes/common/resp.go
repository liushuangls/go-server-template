package common

import (
	"github.com/gin-gonic/gin"
	"github.com/liushuangls/go-server-template/pkg/ecode"
	"go.uber.org/zap"
)

var (
	log *zap.SugaredLogger
)

func SetRespLog(logger *zap.SugaredLogger) {
	log = logger
}

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ErrorResp(c *gin.Context, err error) {
	e := ecode.FromError(err)
	msg := e.Message
	if e.Code == ecode.UnknownCode {
		log.Errorf("ErrorResp receive unknown error: %s", err)
		msg = "Internal Server Error"
	}
	c.JSON(e.HttpCode, Resp{
		Code:    e.Code,
		Message: msg,
		Data:    nil,
	})
	c.Abort()
}

func ParamsErrorResp(c *gin.Context, err error) {
	ErrorResp(c, ecode.NewInvalidParamsErr(err.Error()))
}

func SuccessResp(c *gin.Context, data any) {
	c.JSON(200, Resp{
		Code:    0,
		Message: "",
		Data:    data,
	})
}

func WrapResp(c *gin.Context) func(data any, err error) {
	return func(data any, err error) {
		if err != nil {
			ErrorResp(c, err)
		} else {
			SuccessResp(c, data)
		}
	}
}
