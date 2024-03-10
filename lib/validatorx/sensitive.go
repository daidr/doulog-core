package validatorx

import (
	"github.com/daidr/doulog-core/lib/filter"
	"github.com/go-playground/validator/v10"
)

func NewSensitiveFunc() validator.Func {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		switch t := i.(type) {
		case string:
			return sensitive(t)
		case []string:
			for _, str := range t {
				if v := sensitive(str); !v {
					return false
				}
			}
			return true
		default:
			// 类型不匹配，如果误用此tag也给过
			return true
		}
	}
}

func sensitive(str string) bool {
	v, _ := filter.Text.Validate(str)
	return v
}
