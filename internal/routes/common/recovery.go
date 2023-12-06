package common

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/liushuangls/go-server-template/pkg/ecode"
)

func HandleRecovery(c *gin.Context, err any) {
	ec := &ecode.Error{
		Code:     ecode.PanicCode,
		HttpCode: 500,
		Message:  "Internal Server Error",
	}
	go saveServerLog(c, fmt.Errorf("panic error: %+v", err), ec)
	c.JSON(ec.HttpCode, &Resp{
		Code:    ec.Code,
		Message: ec.Message,
		Data:    nil,
	})
	c.Abort()
}
