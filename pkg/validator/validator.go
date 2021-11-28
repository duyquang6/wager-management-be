package validator

import (
	"fmt"
	"math"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	err := validate.RegisterValidation("monetary-format", validateMonetaryFormat)
	if err != nil {
		panic(err)
	}
}

// GetValidate validate struct fields tag
func GetValidate() *validator.Validate {
	return validate
}

func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "monetary-format":
		return "Invalid monetary format, only accept 2 decimal point"
	case "gt":
		return fmt.Sprintf("This field must be larger than %s", fe.Param())
	case "gte":
		return fmt.Sprintf("This field must be larger or equal %s", fe.Param())
	case "lte":
		return fmt.Sprintf("This field must be lesser or equal %s", fe.Param())
	}
	return fe.Error()
}

// ValidateMonetaryFormat implements validator.Func
func validateMonetaryFormat(fl validator.FieldLevel) bool {
	const epsilon = 1e-8
	floatData := fl.Field().Float()
	return floatData*1e2-math.Floor(floatData*1e2) < epsilon
}
