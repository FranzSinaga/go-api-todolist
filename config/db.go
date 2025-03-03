package config

import (
	"database/sql"
	"fmt"
	_log "log"
	"os"

	"github.com/go-kit/log"
)

func DBConnection(logger log.Logger) *sql.DB {
	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DB")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", username, password, host, dbPort, dbName))
	if err != nil {
		logger.Log("ERROR MESSAGE", err)
		panic(err)
	}
	defer db.Close()
	_log.Printf("DB Connection Successfully")

	rows, err := db.Query("SELECT * FROM todos LIMIT 1")
	if err != nil {
		logger.Log("ERROR MESSAGE", err)
		panic(err)
	}
	rows.Close()

	if db == nil {
		logger.Log("level", "error", "message", "Database connection is nil")
		panic("Database connection failed")
	}
	return db
}
