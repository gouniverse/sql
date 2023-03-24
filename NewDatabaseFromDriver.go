package sql

import (
	"database/sql"
	"errors"
)

func NewDatabaseFromDriver(driverName, dataSourceName string) (*Database, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, errors.New("failed to open DB: " + err.Error())
	}
	return &Database{db: db}, nil
}