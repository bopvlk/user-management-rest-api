package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	err := cv.Validator.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return re.Match([]byte(fl.Field().String()))
	})
	if err != nil {
		return err
	}

	if err := cv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}
