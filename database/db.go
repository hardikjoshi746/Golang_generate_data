package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hardikjoshi746/Golang_generate_data/config"
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

	for i := 0; i < 10; i++ {
		err = DB.Ping()
		if err == nil {
			break
		}
		log.Printf("Waiting for database... (%d/10)\n", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Database not responding after 10 tries:", err)
	}

	DB.SetMaxOpenConns(1000)
	DB.SetMaxIdleConns(1000)
	DB.SetConnMaxLifetime(time.Minute * 5)

	CreateTables()
}
