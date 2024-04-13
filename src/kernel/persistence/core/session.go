package persistence

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Session struct {
	db *sql.DB
}
