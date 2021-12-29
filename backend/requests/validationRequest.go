package requests

import "github.com/go-playground/validator/v10"

type ErrorResponse struct {
	Field string
	Tag   string
	Value string
}

func ValidateStruct(validation ValidateInterface) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(validation)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.Field()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
