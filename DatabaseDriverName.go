package sql

import (
	"database/sql"
	"reflect"
	"strings"
)

// DatabaseDriverName finds the driver name from database
func DatabaseDriverName(db *sql.DB) string {
	dv := reflect.ValueOf(db.Driver())
	driverFullName := dv.Type().String()

	if strings.Contains(driverFullName, "mysql") {
		return "mysql"
	}

	if strings.Contains(driverFullName, "postgres") || strings.Contains(driverFullName, "pq") {
		return "postgres"
	}

	if strings.Contains(driverFullName, "sqlite") {
		return "sqlite"
	}

	if strings.Contains(driverFullName, "mssql") {
		return "mssql"
	}

	return driverFullName
}
