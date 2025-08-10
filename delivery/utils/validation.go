package utils

import "github.com/go-playground/validator/v10"

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			errors[fieldErr.Field()] = getErrorMessage(fieldErr)
		}
	} else {
		// If it's not a validator.ValidationErrors type, return a generic error
		errors["error"] = err.Error()
	}
	return errors
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "gte":
		return "Value must be greater than or equal to " + fe.Param()
	case "lte":
		return "Value must be less than or equal to " + fe.Param()
	default:
		return fe.Error() // default error from validator
	}
}
