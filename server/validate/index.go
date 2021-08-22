package validate

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidRequest(err interface{}) []string {
	var errs []string
	if err == nil {
		return errs
	}
	if validateErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validateErrs {
			sErr := fmt.Errorf("%s 参数错误: %v", e.Namespace(), e.Value())
			errs = append(errs, sErr.Error())
		}
	}
	return errs
}
