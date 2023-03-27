package sql

import (
	"database/sql"
	"testing"

	// _ "github.com/glebarez/go-sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func TestDatabaseDriverName(t *testing.T) {
	conn, err := sql.Open("sqlite3", "test_newdatabase.db")
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	name := DatabaseDriverName(conn)

	if name != "sqlite" {
		t.Fatal(`Error must be "sqlite" but got: `, name)
	}
}
