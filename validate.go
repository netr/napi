package napi

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
)

type CanValidate struct{}

func (v CanValidate) Validate(c *fiber.Ctx, req interface{}) *ValidationError {
	err := c.BodyParser(req)
	if err != nil {
		return &ValidationError{error: err}
	}

	if vErr := validateRequest(req); vErr != nil {
		return &ValidationError{bag: vErr}
	}

	return nil
}

type errorBag map[string]string
type ValidationError struct {
	error error
	bag   errorBag
}

func (ve ValidationError) Error() map[string]string {
	if len(ve.bag) > 0 {
		return ve.bag
	}
	return errorBag{"error": ve.error.Error()}
}

// //////////////////////
// ERROR TRANSLATIONS
// //////////////////////
var errStartsWith string = "{0} must start with '{1}'"

//var validate *validator.Validate

func validateRequest[T any](s T) map[string]string {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)
	addTranslation(validate, trans, "startswith", errStartsWith)

	err := validate.Struct(s)
	if err != nil {
		return translateError(err, trans)
	}
	return nil
}

func addTranslation(validate *validator.Validate, trans ut.Translator, tag string, errMessage string) {
	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, errMessage, false)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		tag := fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}

	_ = validate.RegisterTranslation(tag, trans, registerFn, transFn)
}

func translateError(err error, trans ut.Translator) (errs map[string]string) {
	if err == nil {
		return nil
	}

	errs = make(map[string]string, 0)
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		msg := strings.Replace(translatedErr.Error(), e.Field(), ToCamelCase(e.Field()), -1)
		errs[e.Field()] = msg
		//errs = append(errs, dtos.NewFieldError(e.Field(), e.Value(), msg))
	}
	return errs
}
