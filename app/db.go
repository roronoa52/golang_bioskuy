package app

import (
	"bioskuy/helper"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func GetConnection(config *helper.Config) *sql.DB {

	urlConnection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	db, err := sql.Open(config.DriverName, urlConnection)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}