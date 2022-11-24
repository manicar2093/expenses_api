package validator

import (
	"fmt"
	"net/http"
)

type (
	ValidationError struct {
		Errors interface{} `json:",inline"`
	}
	StructValidable interface {
		ValidateStruct(toValidate interface{}) error
	}
)

func (c *ValidationError) Error() string {
	return fmt.Sprint(c.Errors)
}

func (c *ValidationError) StatusCode() int {
	return http.StatusBadRequest
}
