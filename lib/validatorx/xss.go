package validatorx

import (
	"github.com/feiin/go-xss"
	"github.com/go-playground/validator/v10"
)

func NewXSSFunc() validator.Func {
	return func(fl validator.FieldLevel) bool {
		str, ok := fl.Field().Interface().(string)
		if !ok {
			return true
		}
		return xss.FilterXSS(str, xss.XssOption{WhiteList: xss.GetDefaultWhiteList()}) == str
	}
}
