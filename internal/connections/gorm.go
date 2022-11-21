package connections

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

func GetGormConnection() *gorm.DB {
	postgres.Open(os.Getenv("DATABASE_URL"))
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")))
	if err != nil {
		log.Panicln(err)
	}
	return db
}
