package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"links-service/internal/link"
	"links-service/internal/stat"
	"links-service/internal/user"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file, use default config")
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
		//DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	db.Migrator().AutoMigrate(&link.Link{}, &user.User{}, &stat.Stat{})

}
