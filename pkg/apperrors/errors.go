package apperrors

import (
	"fmt"
	"net/http"
)

type (
	AlreadyExistsError struct {
		Identifier interface{}
		Entity     string
	}
	NotFoundError struct {
		Identifier interface{}
		Message    string
		Entity     string
	}

	MessagedError struct {
		Message string `json:"message,omitempty"`
		Code    int    `json:"-"`
	}
)

func (c AlreadyExistsError) Error() string {
	return fmt.Sprintf("record for entity '%s' with identifier '%s' already exists", c.Entity, c.Identifier)
}

func (c AlreadyExistsError) StatusCode() int {
	return http.StatusBadRequest
}

func (c NotFoundError) Error() string {
	return fmt.Sprintf(
		"record for entity '%s' with identifier '%s' not found: %s",
		c.Entity,
		c.Identifier,
		c.Message,
	)
}

func (c NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

func (c *MessagedError) Error() string {
	return c.Message
}

func (c *MessagedError) StatusCode() int {
	return c.Code
}
