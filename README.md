# SQL

![tests](https://github.com/gouniverse/sql/workflows/tests/badge.svg)

Simple SQL builder

## Example

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

## Installation

```ssh
go get -u github.com/gouniverse/sql
```

## Other
https://github.com/elgris/golang-sql-builder-benchmark
