package db

import (
	"database/sql"
	"log"
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
	makeTable()
	return nil
}

func makeTable() {
	rows, err := Conn.Query("PRAGMA table_info(menu)")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	exists := rows.Next()

	if !exists {
		_, err = Conn.Exec(`CREATE TABLE "menu" (
	"_id"	INTEGER NOT NULL UNIQUE,
	"title"	TEXT,
	"type"	TEXT,
	"ingredients"	TEXT,
	"cook_method"	TEXT,
	"image_list"	TEXT,
	"budget"	REAL,
	"create_at"	TEXT DEFAULT (datetime('now', 'localtime')),
	"update_at"	TEXT DEFAULT (datetime('now', 'localtime')),
	PRIMARY KEY("_id" AUTOINCREMENT)
);`)
		if err != nil {
			log.Println(err)
		}
	}
}

func Close() error {
	return Conn.Close()
}
