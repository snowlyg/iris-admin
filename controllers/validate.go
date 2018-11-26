package controllers

import (
	"IrisYouQiKangApi/tools"
	"gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
	Tools    *tools.Tools
)

func init() {
	validate = validator.New()
	Tools = tools.New()
}

func errorData(errs ...error) string {
	var s string
	for _, err := range errs {
		if err != nil {
			s += err.Error() + "<br/>"
		}
	}
	return s
}
