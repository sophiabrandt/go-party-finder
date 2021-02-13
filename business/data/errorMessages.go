package data

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// TranslateOverrides overrides individual default error messages
// from the validation package with more human-friendly messages.
func TranslateOverride(validate *validator.Validate, trans ut.Translator) {
	validate.RegisterTranslation("required_without", trans, func(ut ut.Translator) error {
		return ut.Add("required_without", "At least 1 Player or 1 Game Master required!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required_without", fe.Field())

		return t
	})
}
