package param_verify

import "gopkg.in/go-playground/validator.v9"

func NameValid(fl validator.FieldLevel) bool {
	if s, ok := fl.Field().Interface().(string); ok {
		if s == "admin" {
			return false
		}
	}
	return true
}
