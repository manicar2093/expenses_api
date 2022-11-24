package validator

import (
	"github.com/gookit/validate"
)

type GooKitValidator struct {
	StopOnError bool
}

const (
	uuidNotValidError = "is not valid for UUID type"
)

func NewGooKitValidator() *GooKitValidator {
	return &GooKitValidator{
		StopOnError: false,
	}
}

func (c *GooKitValidator) ValidateStruct(toValidate interface{}) error {
	v := c.configuratedValidator(toValidate)
	if v.Validate() {
		return nil
	}
	return &ValidationError{Errors: v.Errors}
}

func (c *GooKitValidator) configuratedValidator(toValidate interface{}) *validate.Validation {
	v := validate.Struct(toValidate)
	v.StopOnError = c.StopOnError
	v.AddMessages(map[string]string{
		"uuid":   uuidNotValidError,
		"isUUID": uuidNotValidError,
	})
	return v
}
