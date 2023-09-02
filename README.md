# SQL <a href="https://gitpod.io/#https://github.com/gouniverse/sql" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

![tests](https://github.com/gouniverse/sql/workflows/tests/badge.svg)

An SQL package that wraps the mainstream DB package to allow transparent working
with transactions and simplified SQL builder (with limited functionality).

For a full SQL builder functionality check: https://doug-martin.github.io/goqu


## Installation

```ssh
go get -u github.com/gouniverse/sql
```


## Example Table Creation SQL

```go
import sb "github.com/gouniverse/sql"

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
		Column("updated_at", "datetime", map[string]string{}).
		Column("deleted_at", "datetime", map[string]string{
			"nullable": "yes",
		}).
		Create()
```

## Example Table Drop SQL

```go
import sb "github.com/gouniverse/sql"

sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Drop()
```


## Example Insert SQL

```go
import sb "github.com/gouniverse/sql"
	
sql := sb.NewBuilder(DIALECT_MYSQL).
	Table("cache").
	Insert(map[string]string{
		"ID":         uid.NanoUid(),
		"CacheKey":   token,
		"CacheValue": string(emailJSON),
		"ExpiresAt":  expiresAt.Format("2006-01-02T15:04:05"),
		"CreatedAt":  time.Now().Format("2006-01-02T15:04:05"),
		"UpdatedAt":  time.Now().Format("2006-01-02T15:04:05"),
	})
```

## Example Delete SQL

```go
sql := sb.NewBuilder(DIALECT_MYSQL).
	Table("user").
	Where(sb.Where{
		Column: "id",
		Operator: "==",
		Value: "1",
	}).
	Limit(1).
	Delete()
```

## Initiating Database Instance

1) From existing Go DB instance
```
myDb := sb.NewDatabaseFromDb(sqlDb, DIALECT_MYSQL)
```

3) From driver
```
myDb = sql.NewDatabaseFromDriver("sqlite3", "test.db")
```

## Example SQL Execute

```
myDb := sb.NewDatabaseFromDb(sqlDb, DIALECT_MYSQL)
err := myDb.Exec(sql)
```

## Example Transaction

```go
import _ "modernc.org/sqlite"
import sb "github.com/gouniverse/sql"

myDb = sql.NewDatabaseFromDriver("sqlite3", "test.db")

myDb.BeginTransaction()

err := Database.Exec(sql1)

if err != nil {
	myDb.RollbackTransaction()
	return err
}

err := Database.Exec(sql2)

if err != nil {
	myDb.RollbackTransaction()
	return err
}

myDB.CommitTransaction()

```

## Example Select as Map

Executes a select query and returns map[string]any

```go

mapAny := myDb.SelectToMapAny(sql)

```

Executes a select query and returns map[string]string

```go

mapString := myDb.SelectToMapAny(sql)

```



## Similar

- https://doug-martin.github.io/goqu - Best SQL Builder for Golang
- https://github.com/elgris/golang-sql-builder-benchmark
- https://github.com/es-code/gql

