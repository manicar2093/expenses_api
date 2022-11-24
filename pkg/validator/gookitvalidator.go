package validator

import (
	"github.com/gookit/validate"
)

type GooKitValidator struct {
	StopOnError bool
}

func NewGooKitValidator() *GooKitValidator {
	return &GooKitValidator{
		StopOnError: false,
	}
}

func (c *GooKitValidator) ValidateStruct(toValidate interface{}) error {
	v := validate.Struct(toValidate)
	v.StopOnError = c.StopOnError

	if v.Validate() {
		return nil
	}
	return &ValidationError{Errors: v.Errors}
}
