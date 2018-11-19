package controllers

import (
	"gopkg.in/go-playground/validator.v9"
	"www/pizzaCmsApi/tools"
)

var (
	validate *validator.Validate
	Tools    *tools.Tools
)

func init() {
	validate = validator.New()
	Tools = tools.New()
}

//////////私有方法
/**
 * 返回数据格式不合法的字符串
 * @method ErrorValidate
 */
func errorValidate() string {
	return `{"state": false, "msg": "数据格式不合法"}`
}

func errorData(errs ...error) string {
	var s string
	for _, err := range errs {
		if err != nil {
			s += err.Error() + "<br/>"
			// Tools.Logs(Tools.StructToString(err))
			// Tools.Logs(err.Value)
		}
	}
	return s
}