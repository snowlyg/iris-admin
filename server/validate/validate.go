package validate

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/snowlyg/iris-admin/server/viper_server"
)

// var (
// 	v *Validator
// )

type Validator struct {
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
}

// init
func init() {
	viper_server.Init(getViperConfig())
}

// Instance
func Instance() *Validator {
	switch CONFIG.Locale {
	case "en":
		return newEn()
	case "zh":
		return newZh()
	default:
		return newZh()
	}
}

func newZh() *Validator {
	zh := zh.New()
	uni := ut.New(zh, zh)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("zh")

	validate := validator.New()
	zh_translations.RegisterDefaultTranslations(validate, trans)
	return &Validator{uni: uni, validate: validate, trans: trans}
}

func newEn() *Validator {
	en := en.New()
	uni := ut.New(en, en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)
	return &Validator{uni: uni, validate: validate, trans: trans}
}

type IdBinding struct {
	Id uint `json:"id" uri:"id" validate:"required"`
}

// ShouldBindUri binds the passed struct pointer using the specified binding engine.
func (val *Validator) ShouldBindUri(ctx *gin.Context) (uint, error) {
	var id IdBinding
	if e := ctx.ShouldBindUri(&id); e != nil {
		return 0, e
	}
	if e := val.Translate(id); e != nil {
		return 0, e
	}
	return id.Id, nil
}

func (val *Validator) Translate(s any) error {
	if err := val.validate.Struct(s); err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)
		for _, v := range errs.Translate(val.trans) {
			if v != "" {
				return errors.New(v)
			}
		}
	}
	return nil
}

func (val *Validator) ValidateMap(data map[string]any, rules map[string]any) error {
	mapErrs := val.validate.ValidateMap(data, rules)
	for mk, err := range mapErrs {
		if err != nil {
			// translate all error at once
			errs := err.(validator.ValidationErrors)
			for _, v := range errs.Translate(val.trans) {
				if v != "" {
					return fmt.Errorf("%s%s", mk, v)
				}
			}
		}
	}
	return nil
}
