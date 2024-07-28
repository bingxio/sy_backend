package db

import (
	"database/sql"
	"sy/config"

	_ "github.com/mattn/go-sqlite3"
)

var Conn *sql.DB

func ConnectDB() error {
	db, err := sql.Open("sqlite3", config.C.DbPath)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	Conn = db
	return nil
}

func CloseDB() error {
	return Conn.Close()
}
