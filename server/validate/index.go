package validate

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidRequest(err interface{}) []string {
	errLen := len(err.(validator.ValidationErrors))
	errs := make([]string, errLen)
	if err == nil {
		return errs
	}
	ch := make(chan string, errLen)
	if validateErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validateErrs {
			ve := e
			go func(e validator.FieldError) {
				sErr := fmt.Errorf("%s 参数错误: %v", e.Namespace(), e.Value())
				ch <- sErr.Error()
			}(ve)
		}
	}
	for i := 0; i < errLen; i++ {
		errs[i] = <-ch
	}
	return errs
}
