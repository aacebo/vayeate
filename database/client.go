package database

import (
	"database/sql"
	"os"
	"vayeate/logger"

	_ "github.com/lib/pq"
)

var client *sql.DB
var log = logger.New("database")

func NewClient() *sql.DB {
	if client != nil {
		return client
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION_STRING"))

	if err != nil {
		log.Error(err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Error(err.Error())
	}

	client = db
	return db
}
