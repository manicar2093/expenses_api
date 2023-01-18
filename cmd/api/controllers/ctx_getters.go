package controllers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/manicar2093/goption"
)

func getUserIDAsOptionlaUUID(ctx echo.Context) goption.Optional[uuid.UUID] {
	return goption.Of(uuid.MustParse(ctx.Get("user_id").(string)))
}

func getUserIDAsUUID(ctx echo.Context) uuid.UUID {
	return uuid.MustParse(ctx.Get("user_id").(string))
}
