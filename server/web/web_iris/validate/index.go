package validate

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ValidRequest 请求参数验证
func ValidRequest(err interface{}) []string {
	var errs []string
	if err == nil {
		return errs
	}
	if validateErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validateErrs {
			e := e
			sErr := fmt.Errorf("%s 参数错误: %v", e.Namespace(), e.Value())
			errs = append(errs, sErr.Error())
		}
	}
	return errs
}
