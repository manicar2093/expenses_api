package errors

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func CreateResponseFromError(ctx echo.Context, e error) error {
	errorData := ErrorResponse{Message: e.Error(), Data: e}
	if he, ok := e.(HandleableError); ok {
		return ctx.JSON(he.StatusCode(), errorData)
	}
	return ctx.JSON(http.StatusInternalServerError, errorData)
}
