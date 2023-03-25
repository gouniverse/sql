package sql

import (
	"testing"
	// _ "github.com/glebarez/go-sqlite"
	// _ "github.com/mattn/go-sqlite3"
)

func TestBuilderTableCreateSqlite(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Column("id", "string", map[string]string{
			"primary": "yes",
			"length":  "40",
		}).
		Column("image", "blob", map[string]string{}).
		Column("created_at", "datetime", map[string]string{}).
		Create()

	expected := `CREATE TABLE 'users'("id" TEXT(40) PRIMARY KEY, "image" BLOB, "created_at" DATETIME);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableCreateMysql(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Column("id", "string", map[string]string{
			"primary": "yes",
			"length":  "40",
		}).
		Column("image", "blob", map[string]string{}).
		Column("created_at", "datetime", map[string]string{}).
		Create()

	expected := "CREATE TABLE `users`(`id` VARCHAR(40) PRIMARY KEY, `image` LONGBLOB, `created_at` DATETIME);"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\n but found:\n", sql)
	}
}

func TestBuilderTableCreatePostgres(t *testing.T) {
	sql := NewBuilder(DIALECT_POSTGRES).
		Table("users").
		Column("id", "string", map[string]string{
			"primary": "yes",
			"length":  "40",
		}).
		Column("image", "blob", map[string]string{}).
		Column("created_at", "datetime", map[string]string{}).
		Create()

	expected := `CREATE TABLE "users"("id" TEXT PRIMARY KEY, "image" BYTEA, "created_at" TIMESTAMP);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}
