package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	auth = "root:qwer1234"

	url = "127.0.0.1:3306"
	db  = "shm"
)

var Conn *sql.DB

func ConnectDB() error {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s@tcp(%s)/%s?loc=Local&parseTime=true", auth, url, db),
	)
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
