package sql

import (
	"database/sql"
	"testing"

	// _ "github.com/glebarez/go-sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func TestNewDatabase(t *testing.T) {
	conn, err := sql.Open("sqlite3", "test_newdatabase.db")
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	db := NewDatabase(conn, DIALECT_SQLITE)
	if db == nil {
		t.Fatal("Database MUST NOT BE NIL")
	}
	if db.db == nil {
		t.Fatal("Database db field MUST NOT BE NIL")
	}
	if db.tx != nil {
		t.Fatal("Database tx field MUST BE NIL")
	}
}
