package sql

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/georgysavva/scany/sqlscan"
	"github.com/gouniverse/maputils"
	"github.com/gouniverse/uid"
)

type Database struct {
	db             *sql.DB
	tx             *sql.Tx
	databaseType   string
	sqlLogEnabled  bool
	sqlLog         map[string]string
	sqlDurationLog map[string]time.Duration
	debug          bool
}

func (d *Database) SqlLog() []map[string]string {
	log := []map[string]string{}
	for key, value := range d.sqlLog {
		sqlDuration := d.sqlDurationLog[key]
		log = append(log, map[string]string{
			"sql":  value,
			"time": sqlDuration.String(),
		})
	}
	return log
}

func (d *Database) SqlLogEmpty() {
	d.sqlLog = map[string]string{}
	d.sqlDurationLog = map[string]time.Duration{}
}

func (d *Database) SqlLogLen() int {
	return len(d.sqlLog)
}

func (d *Database) SqlLogEnable(enable bool) {
	d.sqlLogEnabled = enable
}

func (d *Database) SqlLogShrink(leaveLast int) {
	if len(d.sqlLog) <= leaveLast {
		return
	}

	keys := []string{}

	for key, _ := range d.sqlLog {
		keys = append(keys, key)
	}

	tempSqlLog := map[string]string{}
	tempSqlDurationLog := map[string]time.Duration{}
	lastKeys := keys[leaveLast:]
	for _, key := range lastKeys {
		tempSqlLog[key] = d.sqlLog[key]
		tempSqlDurationLog[key] = d.sqlDurationLog[key]
	}

	d.sqlLog = tempSqlLog
	d.sqlDurationLog = tempSqlDurationLog
}

func (d *Database) DebugEnable(debug bool) {
	d.debug = debug
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

func (d *Database) Exec(sqlStr string, args ...any) (sql.Result, error) {
	if d.sqlLogEnabled {
		if d.sqlLog == nil {
			d.sqlLog = map[string]string{}
			d.sqlDurationLog = map[string]time.Duration{}
		}

		sqlID := uid.HumanUid()

		d.sqlLog[sqlID] = sqlStr

		start := time.Now()
		defer func() {
			d.sqlDurationLog[sqlID] = time.Since(start)
		}()
	}

	if d.debug {
		log.Println(sqlStr)
	}

	if d.tx != nil {
		return d.tx.Exec(sqlStr, args...)
	}
	return d.db.Exec(sqlStr, args...)
}

func (d *Database) Query(sqlStr string, args ...any) (*sql.Rows, error) {
	if d.sqlLogEnabled {
		if d.sqlLog == nil {
			d.sqlLog = map[string]string{}
			d.sqlDurationLog = map[string]time.Duration{}
		}

		sqlID := uid.HumanUid()

		d.sqlLog[sqlID] = sqlStr

		start := time.Now()
		defer func() {
			d.sqlDurationLog[sqlID] = time.Since(start)
		}()
	}

	if d.debug {
		log.Println(sqlStr)
	}

	if d.tx != nil {
		return d.tx.Query(sqlStr, args...)
	}
	return d.db.Query(sqlStr, args...)
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
	if d.sqlLogEnabled {
		if d.sqlLog == nil {
			d.sqlLog = map[string]string{}
			d.sqlDurationLog = map[string]time.Duration{}
		}

		sqlID := uid.HumanUid()

		d.sqlLog[sqlID] = sqlStr

		start := time.Now()
		defer func() {
			d.sqlDurationLog[sqlID] = time.Since(start)
		}()
	}

	if d.debug {
		log.Println(sqlStr)
	}

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
