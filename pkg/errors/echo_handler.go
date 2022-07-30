package errors

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateResponseFromError(ctx echo.Context, e error) error {
	errorMap := map[string]interface{}{"error": e.Error()}
	if he, ok := e.(HandleableError); ok {
		return ctx.JSON(he.StatusCode(), errorMap)
	}
	return ctx.JSON(http.StatusInternalServerError, errorMap)
}
