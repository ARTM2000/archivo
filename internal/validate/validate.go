package validate

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("default", defaultTagValidate)
}

type validationError struct {
	FailedField string `json:"field"`
	Message     string `json:"message"`
}

func ValidateStruct[T any](s *T) (errors []*validationError, ok bool) {
	err := validate.Struct(s)
	ok = true
	if err != nil {
		ok = false
		for _, err := range err.(validator.ValidationErrors) {
			var e validationError
			e.FailedField = err.Field()
			if strings.Contains(err.Tag(), "|") {
				tags := strings.Split(err.Tag(), "|")
				last := len(tags) - 1
				infoMsg := strings.Join(tags[:last], ",")
				if len(tags) > 1 {
					infoMsg = strings.Join(tags[:last], ",") + fmt.Sprintf(" or %s", tags[last])
				}
				e.Message = GetValidatorErrorMessage(
					"groupinvalid",
					err.Field(),
					err.Param(),
					infoMsg,
				)
			} else {
				e.Message = GetValidatorErrorMessage(err.Tag(), err.Field(), err.Param())
			}
			errors = append(errors, &e)
		}
	}
	return
}
