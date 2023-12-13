package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func GetDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/p2_ngc4")

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
