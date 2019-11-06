package controllers

import (
	"gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
) 

func init() {
	validate = validator.New()
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
