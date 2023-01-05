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

func (c *GooKitValidator) Validate(i interface{}) error {
	return c.ValidateStruct(i)
}

func (c *GooKitValidator) configuratedValidator(toValidate interface{}) *validate.Validation {
	v := validate.Struct(toValidate) //nolint:varnamelen
	v.StopOnError = c.StopOnError
	v.AddMessages(map[string]string{
		"uuid":     uuidNotValidError,
		"isUUID":   uuidNotValidError,
		"required": "needs to be on request",
	})
	return v
}
