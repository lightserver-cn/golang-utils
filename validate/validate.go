package validate

import (
	"regexp"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zt "github.com/go-playground/validator/v10/translations/zh"
)

var validate *validator.Validate

// 验证
func Validate(v interface{}) (errMaps map[string]string) {
	trans, _ := ut.New(zh.New()).GetTranslator("zh")
	validate = validator.New()
	_ = zt.RegisterDefaultTranslations(validate, trans)

	errMaps = make(map[string]string)
	if err := validate.Struct(v); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if field := err.Field(); field != "" {
				errMaps[field] = err.Translate(trans)
			}
		}
	}

	return errMaps
}

// 正则规则验证
func ValidateRegexTag(tag, pattern string, v *validator.Validate) {
	_ = v.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
		matched, _ := regexp.MatchString(pattern, fl.Field().String())
		return matched
	}, false)
}
