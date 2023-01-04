package connections

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	"github.com/manicar2093/expenses_api/internal/config"
)

func GetGormConnection() *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.Instance.DatabaseURL))
	if err != nil {
		log.Panicln(err)
	}
	return db
}
