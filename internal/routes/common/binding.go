package common

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"

	"github.com/liushuangls/go-server-template/pkg/ecode"
)

type CustomBinder struct {
	b       *echo.DefaultBinder
	v       *validator.Validate
	enTrans ut.Translator
}

func NewCustomBinder() (*CustomBinder, error) {
	cb := &CustomBinder{
		b: new(echo.DefaultBinder),
		v: validator.New(),
	}

	if err := cb.initTrans(); err != nil {
		return nil, err
	}

	return cb, nil
}

func (cb *CustomBinder) initTrans() error {
	enTrans := en.New()
	uni := ut.New(enTrans, enTrans)
	cb.enTrans, _ = uni.GetTranslator("en")
	return entranslations.RegisterDefaultTranslations(cb.v, cb.enTrans)
}

func (cb *CustomBinder) Bind(i any, c echo.Context) (err error) {
	// 1. 先用默认逻辑做绑定（JSON / form / query / param 等）
	if err = cb.b.Bind(i, c); err != nil {
		return ecode.NewInvalidParamsErr("bind error: " + err.Error())
	}

	// 2. 用 go-playground/validator 做结构体校验
	if err = cb.v.Struct(i); err != nil {
		return ecode.NewInvalidParamsErr(cb.translateErr(err))
	}

	return nil
}

func (cb *CustomBinder) translateErr(err error) string {
	if err == nil {
		return ""
	}
	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return err.Error()
	}
	var builder strings.Builder
	for _, e := range errs {
		builder.WriteString(e.Translate(cb.enTrans))
		builder.WriteByte('\n')
	}
	return builder.String()
}
