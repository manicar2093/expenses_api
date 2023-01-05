package middlewares

import "github.com/labstack/echo/v4"

type Middlewares interface {
	LoggedIn(next echo.HandlerFunc) echo.HandlerFunc
}
