package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"links-service/configs"
)

type Db struct {
	*gorm.DB
}

func NewDb(config *configs.Config) *Db {
	conn, err := gorm.Open(postgres.Open(config.Db.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &Db{conn}
}
