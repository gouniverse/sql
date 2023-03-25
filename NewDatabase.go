package sql

import "database/sql"

func NewDatabase(db *sql.DB, databaseType string) *Database {
	return &Database{
		db:           db,
		databaseType: databaseType,
	}
}
