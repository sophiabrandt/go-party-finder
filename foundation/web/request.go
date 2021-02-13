package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-playground/form/v4"
	"github.com/sophiabrandt/go-party-finder/business/data"
	"github.com/sophiabrandt/go-party-finder/business/app/forms"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// translator is a cache of locale and translation information.
var translator *ut.UniversalTranslator

func init() {
	// Instantiate the validator for use.
	validate = validator.New()

	// Instantiate the english locale for the validator library.
	enLocale := en.New()

	// Create a value using English as the fallback locale (first argument).
	// Provide one or more arguments for additional supported locales.
	translator = ut.New(enLocale, enLocale)

	// Register the english error messages for validation errors.
	trans, _ := translator.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	// Translate some error messages to a more human-friendly format.
	data.TranslateOverride(validate, trans)

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Param returns the web call parameter from the request.
func Param(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
//
// If the provided value is a struct then it is checked for validation tags.
func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		// lang controls the language of the error messages.
		lang, _ := translator.GetTranslator("en")

		var fields []FieldError
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Error: verror.Translate(lang),
			}
			fields = append(fields, field)
		}

		return &Error{
			Err:    errors.New("field validation error"),
			Status: http.StatusBadRequest,
			Fields: fields,
		}
	}

	return nil
}

// DecodeForm parses incoming form data and validates it.
func DecodeForm(r *http.Request, val interface{}) (*forms.Form, error) {
	decoder := form.NewDecoder()
	// Use JSON tag names for errors instead of Go struct names.
	decoder.SetTagName("json")

	err := r.ParseForm()
	form := forms.New(r.PostForm)

	if err != nil {
		return form, NewRequestError(err, http.StatusBadRequest)
	}
	if err := decoder.Decode(&val, r.Form); err != nil {
		return form, NewRequestError(err, http.StatusBadRequest)
	}

	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return form, err
		}

		// lang controls the language of the error messages.
		lang, _ := translator.GetTranslator("en")

		// add the validation errors as form errors
		for _, verror := range verrors {
			form.Errors.Add(verror.Field(), verror.Translate(lang))
		}
		return form, err
	}

	return form, nil
}
