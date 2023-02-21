package validator

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {

	err := cv.Validator.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		for _, test := range []string{`.{7,}`, `[\p{Lu}]`, `[\p{Ll}]`, "[0-9]", `[^\\d\\w]`} {
			t, _ := regexp.MatchString(test, fl.Field().String())
			if !t {
				return false
			}
		}
		return true
	})
	if err != nil {
		return err
	}

	if err := cv.Validator.Struct(i); err != nil {
		return fmt.Errorf("password must containe at least 7 letters,  1 number, 1 upper case, 1 special character.  err: %v", err)
	}
	return nil
}
