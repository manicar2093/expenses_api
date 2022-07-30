package repos

import (
	"fmt"
	"net/http"
)

type (
	AlreadyExistsError struct {
		Identifier string
		Entity     string
	}
)

func (c AlreadyExistsError) Error() string {
	return fmt.Sprintf("record for entity '%s' with identifier '%s' already exists", c.Entity, c.Identifier)
}

func (c AlreadyExistsError) StatusCode() int {
	return http.StatusBadRequest
}
