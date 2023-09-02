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
		Column("price_default", "decimal", map[string]string{}).
		Column("price_custom", "decimal", map[string]string{
			"length":   "12",
			"decimals": "10",
		}).
		Column("created_at", "datetime", map[string]string{}).
		Column("deleted_at", "datetime", map[string]string{
			"nullable": "yes",
		}).
		Create()

	expected := `CREATE TABLE "users"("id" TEXT(40) PRIMARY KEY NOT NULL, "image" BLOB NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "created_at" DATETIME NOT NULL, "deleted_at" DATETIME);`
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
		Column("price_default", "decimal", map[string]string{}).
		Column("price_custom", "decimal", map[string]string{
			"length":   "12",
			"decimals": "10",
		}).
		Column("created_at", "datetime", map[string]string{}).
		Column("deleted_at", "datetime", map[string]string{
			"nullable": "yes",
		}).
		Create()

	expected := "CREATE TABLE `users`(`id` VARCHAR(40) PRIMARY KEY NOT NULL, `image` LONGBLOB NOT NULL, `price_default` DECIMAL(10,2) NOT NULL, `price_custom` DECIMAL(12,10) NOT NULL, `created_at` DATETIME NOT NULL, `deleted_at` DATETIME);"
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
		Column("price_default", "decimal", map[string]string{}).
		Column("price_custom", "decimal", map[string]string{
			"length":   "12",
			"decimals": "10",
		}).
		Column("created_at", "datetime", map[string]string{}).
		Column("deleted_at", "datetime", map[string]string{
			"nullable": "yes",
		}).
		Create()

	expected := `CREATE TABLE "users"("id" TEXT PRIMARY KEY NOT NULL, "image" BYTEA NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "created_at" TIMESTAMP NOT NULL, "deleted_at" TIMESTAMP);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDropMysql(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Drop()

	expected := "DROP TABLE `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDropPostgres(t *testing.T) {
	sql := NewBuilder(DIALECT_POSTGRES).
		Table("users").
		Drop()

	expected := `DROP TABLE "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDropSqlite(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Drop()

	expected := `DROP TABLE "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteMysql(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Delete()

	expected := "DELETE FROM `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteMysqlExtended(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Where(Where{
			Column:   "FirstName",
			Operator: "==",
			Value:    "Tom",
		}).
		Where(Where{
			Column:   "FirstName",
			Operator: "==",
			Value:    "Sam",
			Type:     "OR",
		}).
		Limit(12).
		Offset(34).
		Delete()

	expected := "DELETE FROM `users` WHERE `FirstName` = \"Tom\" OR `FirstName` = \"Sam\" LIMIT 12 OFFSET 34;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteSqlite(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Delete()

	expected := `DELETE FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteSqliteExtended(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Where(Where{
			Column:   "FirstName",
			Operator: "==",
			Value:    "Tom",
		}).
		Where(Where{
			Column:   "FirstName",
			Operator: "==",
			Value:    "Sam",
			Type:     "OR",
		}).
		Delete()

	expected := `DELETE FROM "users" WHERE "FirstName" = 'Tom' OR "FirstName" = 'Sam';`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}
