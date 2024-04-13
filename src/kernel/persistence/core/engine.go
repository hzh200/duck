package persistence

import "database/sql"

type Engine struct {
	db *sql.DB
	dialect Dialect
}
