package web_gin

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/snowlyg/iris-admin/server/web"
)

func registerValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("dev-required", validateDevRequired)
	}
}

var validateDevRequired validator.Func = func(fl validator.FieldLevel) bool {
	if web.CONFIG.System.Level == "release" {
		return fl.Field().String() != ""
	}
	return true
}
