package common

import (
	"net/http"

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
		log.Errorf("NewResp receive unknown error: %s", ec)
	}
	return httpCode, &Resp{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

func ErrorResp(c *gin.Context, err error) {
	c.JSON(NewResp(nil, err))
	c.Abort()
}

func ParamsErrorResp(c *gin.Context, err error) {
	ErrorResp(c, ecode.NewInvalidParamsErr(translateErr(err)))
}

func SuccessResp(c *gin.Context, data any) {
	c.JSON(NewResp(data, nil))
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
