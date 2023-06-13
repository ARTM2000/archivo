package validate

import (
	"reflect"

	"log"

	"github.com/go-playground/validator/v10"
)

func defaultTagValidate(fl validator.FieldLevel) bool {
	// get default value
	defaultTagValue := fl.Param()
	if defaultTagValue != "" {
		// if there is no default value specified, ignore tag
		panic("default value for validation tag 'default' not found")
	}
	log.Default().Println("h > ", fl.Field(), reflect.TypeOf(fl.Field().Interface()), defaultTagValue)
	return true

	// fieldType := reflect.TypeOf(fl.Field().Interface());

}
