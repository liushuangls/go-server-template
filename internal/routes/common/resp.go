package common

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"github.com/liushuangls/go-server-template/internal/data"
	"github.com/liushuangls/go-server-template/internal/data/ent"
	"github.com/liushuangls/go-server-template/internal/data/ent/serverlog"
	"github.com/liushuangls/go-server-template/pkg/ecode"
	"github.com/liushuangls/go-server-template/pkg/xstrings"
)

var (
	logIncludeCode = []int{ecode.InternalServerErr.Code}
)

var (
	serverLog *data.ServerLogRepo
)

func SetServerLogRepo(repo *data.ServerLogRepo) {
	serverLog = repo
}

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    any    `json:"data"`
}

func NewResp(c *gin.Context, data interface{}, err error) (int, *Resp) {
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
		//slog.Error("NewResp receive unknown error", "err", ec)
	}
	if err != nil && ec != nil {
		go saveServerLog(c, err, ec)
	}
	return httpCode, &Resp{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

func ErrorResp(c *gin.Context, err error) {
	c.JSON(NewResp(c, nil, err))
	c.Abort()
}

func ParamsErrorResp(c *gin.Context, err error) {
	var errs *ecode.Error
	if errors.As(err, &errs) {
		ErrorResp(c, errs)
		return
	}
	ErrorResp(c, ecode.NewInvalidParamsErr(translateErr(err)))
}

func SuccessResp(c *gin.Context, data any) {
	c.JSON(NewResp(c, data, nil))
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

func saveServerLog(c *gin.Context, err error, ec *ecode.Error) {
	var (
		uid         int
		body        []byte
		contentType = c.Request.Header.Get("Content-Type")
		level       = serverlog.LevelERROR
		clientIP    string
	)
	switch ec.Code {
	case ecode.PanicCode:
		level = serverlog.LevelPANIC
	case ecode.UnknownCode:
		level = serverlog.LevelERROR
	default:
		if lo.Contains(logIncludeCode, ec.Code) {
			level = serverlog.LevelERROR
		} else {
			level = serverlog.LevelWARN
		}
	}
	if user := GetCurrentUserInfo(c); user != nil {
		uid = user.ID
	}
	if ipInfo := GetIPInfo(c); ipInfo != nil {
		clientIP = ipInfo.IP
	}
	if strings.Contains(contentType, "application/json") {
		// 只有route层调用了 c.ShouldBindBodyWith 后这里才能读到
		if cb, ok := c.Get(gin.BodyBytesKey); ok {
			if cbb, ok := cb.([]byte); ok {
				body = cbb
			}
		}
	}
	if len(body) > 65000 {
		body = []byte("body greater than 65535")
	}
	aLog := &ent.ServerLog{
		UserID:     uid,
		IP:         clientIP,
		Method:     c.Request.Method,
		Path:       c.Request.URL.Path,
		Query:      c.Request.URL.Query().Encode(),
		ErrMsg:     err.Error(),
		RespErrMsg: ec.Log(),
		Body:       xstrings.BytesToString(body),
		Level:      level,
		Code:       ec.Code,
	}
	if _, err := serverLog.Create(context.Background(), aLog); err != nil {
		slog.Error("routes.common.saveServerLog", "error", err)
	}
}
