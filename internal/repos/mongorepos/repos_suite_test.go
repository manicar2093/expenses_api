package mongorepos_test

import (
	"github.com/manicar2093/expenses_api/internal/connections"
)

var (
	conn = connections.GetMongoConn()
)
