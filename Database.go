package sql

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/georgysavva/scany/sqlscan"
	"github.com/gouniverse/maputils"
)

type Database struct {
	db           *sql.DB
	tx           *sql.Tx
	databaseType string
}

func (d *Database) Type() string {
	return d.databaseType
}

func (d *Database) Open() (err error) {
	d.db, err = sql.Open("sqlite3", "test.db")
	return err
}

func (d *Database) Close() (err error) {
	return d.db.Close()
}

func (d *Database) DB() *sql.DB {
	return d.db
}

func (d *Database) BeginTransaction() (err error) {
	if d.tx != nil {
		return errors.New("transaction already in progress")
	}

	tx, err := d.db.Begin()
	if err != nil {
		return errors.New("failed to begin transaction: " + err.Error())
	}
	d.tx = tx

	return err
}

func (d *Database) BeginTransactionWithContext(ctx context.Context, opts *sql.TxOptions) (err error) {
	if d.tx != nil {
		return errors.New("transaction already in progress")
	}

	tx, err := d.db.BeginTx(ctx, opts)

	if err != nil {
		return errors.New("failed to begin transaction: " + err.Error())
	}

	d.tx = tx

	return nil
}

func (d *Database) ExecInTransaction(fn func(d *Database) error) (err error) {
	err = d.BeginTransaction()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err := d.RollbackTransaction()
			if err != nil {
				log.Println("sqldb rollback error: " + err.Error())
			}
		}
	}()

	err = fn(&Database{db: d.db, tx: d.tx})

	if err == nil {
		err = d.CommitTransaction()
	}

	return
}

func (d *Database) Exec(sql string, args ...any) (sql.Result, error) {
	if d.tx != nil {
		return d.tx.Exec(sql, args...)
	}
	return d.db.Exec(sql, args...)
}

func (d *Database) Query(sql string, args ...any) (*sql.Rows, error) {
	if d.tx != nil {
		return d.tx.Query(sql, args...)
	}
	return d.db.Query(sql, args...)
}

func (d *Database) CommitTransaction() (err error) {
	if d.tx == nil {
		return errors.New("no transaction in progress")
	}

	err = d.tx.Commit()

	if err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	d.tx = nil // empty transaction

	return err
}

func (d *Database) RollbackTransaction() (err error) {
	if d.tx == nil {
		return errors.New("no transaction in progress")
	}

	err = d.tx.Rollback()

	if err != nil {
		return errors.New("failed to rollback transaction: " + err.Error())
	}

	d.tx = nil // empty transaction

	return err
}

func (d *Database) SelectToMapAny(sqlStr string, args ...any) ([]map[string]any, error) {
	listMap := []map[string]any{}

	err := sqlscan.Select(context.Background(), d.db, &listMap, sqlStr)
	if err != nil {
		if sqlscan.NotFound(err) {
			return []map[string]any{}, nil
		}

		return []map[string]any{}, err
	}

	return listMap, nil
}

func (d *Database) SelectToMapString(sqlStr string, args ...any) ([]map[string]string, error) {
	listMapAny, err := d.SelectToMapAny(sqlStr, args...)

	if err != nil {
		return []map[string]string{}, err
	}

	listMapString := []map[string]string{}

	for i := 0; i < len(listMapAny); i++ {
		mapString := maputils.MapStringAnyToMapStringString(listMapAny[i])
		listMapString = append(listMapString, mapString)
	}

	return listMapString, nil
}
