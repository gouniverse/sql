# SQL

![tests](https://github.com/gouniverse/sql/workflows/tests/badge.svg)

Simple SQL builder


## Installation

```ssh
go get -u github.com/gouniverse/sql
```


## Example Table Creation

```go
import sb "github.com/gouniverse/sql"



sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Column("id", "string", map[string]string{
			"primary": "yes",
			"length":  "40",
		}).
		Column("image", "blob", map[string]string{}).
		Column("created_at", "datetime", map[string]string{}).
		Create()

myDb := sb.NewDatabaseFromDb(sqlDb, DIALECT_MYSQL)
err := myDb.Exec(sql)
```

## Example Transaction

```go
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



## Example of the Otdated Builder (do not use)

This builder was version 1 and is now outdated

```go
import sb "github.com/gouniverse/sql"
	
sql := sb.NewSqlite().Table("cache").Insert(map[string]string{
		"ID":         uid.NanoUid(),
		"CacheKey":   token,
		"CacheValue": string(emailJSON),
		"ExpiresAt":  expiresAt.Format("2006-01-02T15:04:05"),
		"CreatedAt":  time.Now().Format("2006-01-02T15:04:05"),
		"UpdatedAt":  time.Now().Format("2006-01-02T15:04:05"),
})
```



## Similar

- https://doug-martin.github.io/goqu - Best SQL Builder
- https://github.com/elgris/golang-sql-builder-benchmark
- https://github.com/es-code/gql

