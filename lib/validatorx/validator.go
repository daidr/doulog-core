package validatorx

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func MustRegister(tag string, fn validator.Func) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("expect *validator.Validate")
	}

	if err := v.RegisterValidation(tag, fn); err != nil {
		panic(err)
	}
}
