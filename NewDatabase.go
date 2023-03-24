package sql

import "database/sql"

func NewDatabase(db *sql.DB) *Database {
	return &Database{db: db}
}
