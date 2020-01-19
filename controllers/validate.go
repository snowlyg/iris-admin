package controllers

import (
	"fmt"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 表单验证
// https://github.com/go-playground/validator
var (
	uni           *ut.UniversalTranslator
	Validate      *validator.Validate
	ValidateTrans ut.Translator
)

func init() {
	zh2 := zh.New()
	uni = ut.New(zh2, zh2)

	ValidateTrans, _ = uni.GetTranslator("zh")

	Validate = validator.New()
	if err := zh_translations.RegisterDefaultTranslations(Validate, ValidateTrans); err != nil {
		fmt.Println(fmt.Sprintf("RegisterDefaultTranslations %v", err))
	}
}
