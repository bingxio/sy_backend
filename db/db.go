package db

import (
	"database/sql"
	"sy_backend/config"

	_ "github.com/mattn/go-sqlite3"
)

var Conn *sql.DB

func Open() error {
	db, err := sql.Open("sqlite3", config.Conf.DbPath)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	Conn = db
	return nil
}

func Close() error {
	return Conn.Close()
}
