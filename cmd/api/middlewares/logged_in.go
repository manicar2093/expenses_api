package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/internal/auth"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
)

// LoggedIn checks if there is a Authorization header. If it is there added it to context
// otherwise return an error
func (c *EchoMiddlewares) LoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authorization := ctx.Request().Header.Get("Authorization")
		if authorization == "" {
			return ctx.JSON(http.StatusBadRequest, errors.New("no token in request"))
		}
		var accessToken auth.AccessToken
		err := c.tokenValidable.ValidateToken(ctx.Request().Context(), strings.Split(authorization, " ")[1], &accessToken)
		if err != nil {
			return apperrors.CreateResponseFromError(ctx, err)
		}
		ctx.Set("user_id", accessToken.UserID)
		return nil
	}
}
