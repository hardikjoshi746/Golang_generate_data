package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"example.com/rest/config"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	cfg := config.Load()

	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	DB, err = sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("cannot connect to database:", err) // if the connection fails the complete app stops
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("database not responding:", err) // if the connection fails the complete app stops
	}

	DB.SetMaxOpenConns(1000)
	DB.SetMaxIdleConns(1000)
	DB.SetConnMaxLifetime(time.Minute * 5)

	CreateTables()
}
