package sql

import (
	"strings"

	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

type Builder struct {
	Dialect      string
	TableName    string
	sql          map[string]any
	sqlColumns   []map[string]any
	sqlTableName string
}

func (b *Builder) Table(tableName string) *Builder {
	b.sqlTableName = tableName
	return b
}

func (b *Builder) Column(columnName string, columnType string, opts map[string]string) *Builder {
	sqlColumn := map[string]any{
		"column_name":    columnName,
		"column_type":    columnType,
		"column_options": opts,
	}
	b.sqlColumns = append(b.sqlColumns, sqlColumn)
	return b
}

/** The column method specifies the desired columns in a table.
 * <code>
 * // Creating table USERS with two columns USER_ID, and USER_NAME
 * $database->table("USERS")
 *     ->column("USER_ID","INTEGER")
 *     ->column("USER_NAME","STRING")
 *     ->create();
 * </code>
 * @param String the name of the column
 * @param String the type of the column (STRING, INTEGER, FLOAT, TEXT, BLOB)
 * @param String the attributes of the column (NOT NULL PRIMARY KEY AUTO_INCREMENT)
 * @return SqlDb an instance of this database
 * @access public
 */
//  function column($column_name, $column_type = null, $column_properties = null)
//  {
// 	 if (isset($this->sql["table"]) == false) {
// 		 trigger_error('ERROR: In class <b>' . get_class($this) . '</b> in method <b>column($column,$details)</b>: Trying to attach column to non-specified table!', E_USER_ERROR);
// 	 }

// 	 $current_table = (count($this->sql["table"]) - 1);
// 	 $current_table_name = $this->sql["table"][$current_table];

// 	 if (isset($this->sql["columns"]) == false) {
// 		 $this->sql["columns"] = array();
// 	 }

// 	 if (isset($this->sql["columns"][$current_table_name]) == false) {
// 		 $this->sql["columns"][$current_table_name] = array();
// 	 }

// 	 $this->sql["columns"][$current_table_name][] = array($column_name, $column_type, $column_properties);
// 	 return $this;
//  }

/**
 * The create method creates new database or table.
 * If the database or table can not be created it will return false.
 * False will be returned if the database or table already exist.
 * <code>
 * // Creating a new database
 * $database->create();
 *
 * // Creating a new table
 * $database->table("STATES")
 *     ->column("STATE_NAME","STRING")
 *     ->create();
 * </code>
 * @return boolean true, on success, false, otherwise
 * @access public
 */

func (b *Builder) Create() string {
	sql := ""
	if b.Dialect == DIALECT_MYSQL {
		sql = "CREATE TABLE `" + b.sqlTableName + "`(" + b.columnsToSQL(b.sqlColumns) + ");"
	}
	if b.Dialect == DIALECT_POSTGRES {
		sql = `CREATE TABLE "` + b.sqlTableName + `"(` + b.columnsToSQL(b.sqlColumns) + `);`
	}
	if b.Dialect == DIALECT_SQLITE {
		sql = "CREATE TABLE '" + b.sqlTableName + "'(" + b.columnsToSQL(b.sqlColumns) + ");"
	}
	return sql
}

func (b *Builder) CreateIfNotExists() string {
	sql := ""
	if b.Dialect == DIALECT_MYSQL {
		sql = "CREATE TABLE IF NOT EXISTS `" + b.sqlTableName + "`(" + b.columnsToSQL(b.sqlColumns) + ");"
	}
	if b.Dialect == DIALECT_POSTGRES {
		sql = `CREATE TABLE IF NOT EXISTS "` + b.sqlTableName + `"(` + b.columnsToSQL(b.sqlColumns) + `);`
	}
	if b.Dialect == DIALECT_SQLITE {
		sql = "CREATE TABLE IF NOT EXISTS '" + b.sqlTableName + "'(" + b.columnsToSQL(b.sqlColumns) + ");"
	}
	return sql
}

// columnsToSQL converts the columns statements to SQL.
func (b *Builder) columnsToSQL(columns []map[string]any) string {
	columnSqls := []string{}

	for i := 0; i < len(columns); i++ {
		column := columns[i]
		columnName := utils.ToString(column["column_name"])
		columnType := utils.ToString(column["column_type"])
		columnOptions := column["column_options"].(map[string]string)
		columnLength := lo.ValueOr(columnOptions, "length", "")
		columnDecimals := lo.ValueOr(columnOptions, "decimals", "")
		columnAuto := lo.ValueOr(columnOptions, "auto", "no")
		columnPrimary := lo.ValueOr(columnOptions, "primary", "no")
		columnNullable := lo.ValueOr(columnOptions, "nullable", "no")

		columnSql := lo.IfF(b.Dialect == DIALECT_MYSQL, func() string {
			columnType := lo.
				IfF(columnType == "string", func() string {
					columnLength = lo.Ternary(columnLength == "", "255", columnLength)
					return "VARCHAR"
				}).
				ElseIfF(columnType == "integer", func() string {
					return "BIGINT"
				}).
				ElseIfF(columnType == "float", func() string {
					return "DOUBLE"
				}).
				ElseIfF(columnType == "text", func() string {
					return "LONGTEXT"
				}).
				ElseIfF(columnType == "blob", func() string {
					return "LONGBLOB"
				}).
				ElseIfF(columnType == "date", func() string {
					return "DATE"
				}).
				ElseIfF(columnType == "datetime", func() string {
					return "DATETIME"
				}).
				ElseIfF(columnType == "decimal", func() string {
					return "DECIMAL"
				}).
				Else(columnType)

			sql := "`" + columnName + "` " + columnType

			// Column length
			if columnType == "DECIMAL" {
				if columnLength == "" {
					columnLength = "10"
				}
				if columnDecimals == "" {
					columnDecimals = "2"
				}
				sql += "(" + columnLength + "," + columnDecimals + ")"

			} else if columnLength != "" {
				sql += "(" + columnLength + ")"
			}

			// Auto increment
			if columnAuto == "yes" {
				sql += " AUTO_INCREMENT"
			}

			// Primary key
			if columnPrimary == "yes" {
				sql += " PRIMARY KEY"
			}

			// Non Nullable / Required
			if columnNullable != "yes" {
				sql += " NOT NULL"
			}
			return sql
		}).ElseIfF(b.Dialect == DIALECT_POSTGRES, func() string {
			columnType := lo.
				IfF(columnType == "string", func() string {
					return "TEXT"
				}).
				ElseIfF(columnType == "integer", func() string {
					return "INTEGER"
				}).
				ElseIfF(columnType == "float", func() string {
					return "REAL"
				}).
				ElseIfF(columnType == "text", func() string {
					return "TEXT"
				}).
				ElseIfF(columnType == "blob", func() string {
					return "BYTEA"
				}).
				ElseIfF(columnType == "date", func() string {
					return "DATE"
				}).
				ElseIfF(columnType == "datetime", func() string {
					return "TIMESTAMP"
				}).
				ElseIfF(columnType == "decimal", func() string {
					return "DECIMAL"
				}).
				Else(columnType)

			sql := `"` + columnName + `" ` + columnType + ``

			// Column length
			if columnType == "DECIMAL" {
				if columnLength == "" {
					columnLength = "10"
				}
				if columnDecimals == "" {
					columnDecimals = "2"
				}
				sql += "(" + columnLength + "," + columnDecimals + ")"

			} else if columnLength != "" && columnType != "TEXT" {
				sql += "(" + columnLength + ")"
			}

			// Auto increment
			if columnAuto == "yes" {
				sql += " SERIAL"
			}

			// Primary key
			if columnPrimary == "yes" {
				sql += " PRIMARY KEY"
			}

			// Non Nullable / Required
			if columnNullable != "yes" {
				sql += " NOT NULL"
			}
			return sql
		}).ElseIfF(b.Dialect == DIALECT_SQLITE, func() string {
			columnType := lo.
				IfF(columnType == "string", func() string {
					return "TEXT"
				}).
				ElseIfF(columnType == "integer", func() string {
					return "INTEGER"
				}).
				ElseIfF(columnType == "float", func() string {
					return "REAL"
				}).
				ElseIfF(columnType == "text", func() string {
					return "TEXT"
				}).
				ElseIfF(columnType == "blob", func() string {
					return "BLOB"
				}).
				ElseIfF(columnType == "date", func() string {
					return "DATE"
				}).
				ElseIfF(columnType == "datetime", func() string {
					return "DATETIME"
				}).
				ElseIfF(columnType == "decimal", func() string {
					return "DECIMAL"
				}).
				Else(columnType)

			sql := `"` + columnName + `" ` + columnType + ``

			// Column length
			if columnType == "DECIMAL" {
				if columnLength == "" {
					columnLength = "10"
				}
				if columnDecimals == "" {
					columnDecimals = "2"
				}
				sql += "(" + columnLength + "," + columnDecimals + ")"

			} else if columnLength != "" {
				sql += "(" + columnLength + ")"
			}

			// Auto increment
			if columnAuto == "yes" {
				sql += " AUTOINCREMENT"
			}

			// Primary key
			if columnPrimary == "yes" {
				sql += " PRIMARY KEY"
			}

			// Non Nullable / Required
			if columnNullable != "yes" {
				sql += " NOT NULL"
			}
			return sql
		}).ElseF(func() string {
			return "not supported"
		})

		columnSqls = append(columnSqls, columnSql)
	}

	return strings.Join(columnSqls, ", ")
}
