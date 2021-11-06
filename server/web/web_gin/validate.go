package web_gin

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func registerValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("dev-required", validateDevRequired)
	}
}

var validateDevRequired validator.Func = func(fl validator.FieldLevel) bool {
	if CONFIG.System.Level != "release" {
		return true
	}
	return fl.Field().String() != ""
}
