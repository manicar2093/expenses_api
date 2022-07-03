package connections

import (
	"os"

	"github.com/go-rel/postgres"
	"github.com/go-rel/rel"

	_ "github.com/lib/pq"
)

func GetRelConnection() rel.Repository {
	adapter, err := postgres.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	return rel.New(adapter)
}
