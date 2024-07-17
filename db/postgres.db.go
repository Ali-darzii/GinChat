package db

import (
	"GinChat/entity"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func ConnectPostgres() *gorm.DB {
	end := godotenv.Load()
	if end != nil {
		panic("Failed to load .env file")
	}
	dbUser := os.Getenv("PostgresUser")
	dbPASS := os.Getenv("PostgresPass")
	dbHost := os.Getenv("PostgresHost")
	dbPort := os.Getenv("PostgresPort")
	dbName := os.Getenv("PostgresName")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
		dbHost, dbUser, dbPASS, dbName, dbPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect Postgres database")
	}

	err = db.AutoMigrate(
		//auth
		&entity.User{},
		&entity.Phone{},
		&entity.UserLogins{},
		&entity.UserDevice{},
		&entity.UserIP{},
		//chat
		&entity.PrivateRoom{},
		&entity.PrivateMessageRoom{},
	)
	if err != nil {
		panic("Failed: Unable to migrate your postgres database")
	}
	return db
}

func ClosePostgres(db *gorm.DB) {
	dbPsql, err := db.DB()
	if err != nil {
		panic("Failed: postgres database connection !")
	}
	err = dbPsql.Close()

	if err != nil {
		panic("Failed: unable to close postgres connection database !")
	}
}
