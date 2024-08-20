package common

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/liushuangls/go-server-template/pkg/ecode"
)

func HandleRecovery(c *gin.Context, err any) {
	ec := &ecode.Error{
		Code:     ecode.PanicCode,
		HttpCode: 500,
		Message:  "Internal Server Error",
	}
	resp := &Resp{
		Code:    ec.Code,
		Message: ec.Message,
		Data:    nil,
	}

	path := c.Request.URL.Path
	if strings.Contains(path, "/chat/completions") {
		c.SSEvent("", resp)
	} else {
		c.JSON(ec.HttpCode, resp)
	}

	c.Abort()
}
