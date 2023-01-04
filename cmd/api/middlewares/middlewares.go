package middlewares

import "github.com/manicar2093/expenses_api/internal/auth"

type EchoMiddlewares struct {
	tokenValidable auth.TokenValidable
}

func NewEchoMiddlewares(tokenValidable auth.TokenValidable) *EchoMiddlewares {
	return &EchoMiddlewares{tokenValidable: tokenValidable}
}
