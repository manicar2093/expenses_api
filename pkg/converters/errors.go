package converters

import "net/http"

type IDNotValidIDError struct{}

func (c IDNotValidIDError) Error() string {
	return "given id it is not valid"
}

func (c IDNotValidIDError) StatusCode() int {
	return http.StatusBadRequest
}
