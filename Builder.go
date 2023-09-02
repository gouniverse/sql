package sql

import (
	"strconv"
	"strings"

	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

type Where struct {
	Raw      string
	Column   string
	Operator string
	Type     string
	Value    string
}

type OrderBy struct {
	Column    string
	Direction string
}

type Builder struct {
	Dialect string
	//TableName    string
	sql          map[string]any
	sqlColumns   []map[string]any
	sqlLimit     int64
	sqlOffset    int64
	sqlOrderBy   []OrderBy
	sqlTableName string
	sqlWhere     []Where
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
		sql = "CREATE TABLE " + b.quoteTable(b.sqlTableName) + "(" + b.columnsToSQL(b.sqlColumns) + ");"
	}
	if b.Dialect == DIALECT_POSTGRES {
		sql = `CREATE TABLE ` + b.quoteTable(b.sqlTableName) + `(` + b.columnsToSQL(b.sqlColumns) + `);`
	}
	if b.Dialect == DIALECT_SQLITE {
		sql = "CREATE TABLE " + b.quoteTable(b.sqlTableName) + "(" + b.columnsToSQL(b.sqlColumns) + ");"
	}
	return sql
}

func (b *Builder) CreateIfNotExists() string {
	sql := ""
	if b.Dialect == DIALECT_MYSQL {
		sql = "CREATE TABLE IF NOT EXISTS " + b.quoteTable(b.sqlTableName) + "(" + b.columnsToSQL(b.sqlColumns) + ");"
	}
	if b.Dialect == DIALECT_POSTGRES {
		sql = `CREATE TABLE IF NOT EXISTS ` + b.quoteTable(b.sqlTableName) + `(` + b.columnsToSQL(b.sqlColumns) + `);`
	}
	if b.Dialect == DIALECT_SQLITE {
		sql = "CREATE TABLE IF NOT EXISTS " + b.quoteTable(b.sqlTableName) + "'(" + b.columnsToSQL(b.sqlColumns) + ");"
	}
	return sql
}

/**
 * The delete method deletes a row in a table. For deleting a database
 * or table use the drop method.
 * <code>
 * // Deleting a row
 * $database->table("STATES")->where("STATE_NAME","=","Alabama")->delete();
 * </code>
 * @return boolean true, on success, false, otherwise
 * @access public
 */
// Drop deletes a table
func (b *Builder) Delete() string {
	if b.sqlTableName == "" {
		panic("In method Delete() no table specified to delete from!")
	}

	where := ""
	if len(b.sqlWhere) > 0 {
		where = b.whereToSql(b.sqlWhere)
	}

	orderBy := ""
	if len(b.sqlOrderBy) > 0 {
		orderBy = b.orderByToSql(b.sqlOrderBy)
	}

	limit := ""
	if b.sqlLimit > 0 {
		limit = " LIMIT " + strconv.FormatInt(b.sqlLimit, 10)
	}

	offset := ""
	if b.sqlOffset > 0 {
		offset = " OFFSET " + strconv.FormatInt(b.sqlOffset, 10)
	}

	sql := ""
	if b.Dialect == DIALECT_MYSQL {
		sql = "DELETE FROM " + b.quoteTable(b.sqlTableName) + where + orderBy + limit + offset + ";"
	}
	if b.Dialect == DIALECT_POSTGRES {
		sql = `DELETE FROM ` + b.quoteTable(b.sqlTableName) + where + orderBy + limit + offset + `;`
	}
	if b.Dialect == DIALECT_SQLITE {
		sql = "DELETE FROM " + b.quoteTable(b.sqlTableName) + where + orderBy + limit + offset + ";"
	}
	return sql
}

// Drop deletes a table
func (b *Builder) Drop() string {
	sql := ""
	if b.Dialect == DIALECT_MYSQL {
		sql = "DROP TABLE " + b.quoteTable(b.sqlTableName) + ";"
	}
	if b.Dialect == DIALECT_POSTGRES {
		sql = `DROP TABLE ` + b.quoteTable(b.sqlTableName) + `;`
	}
	if b.Dialect == DIALECT_SQLITE {
		sql = "DROP TABLE " + b.quoteTable(b.sqlTableName) + ";"
	}
	return sql
}

func (b *Builder) Limit(limit int64) *Builder {
	b.sqlLimit = limit
	return b
}

func (b *Builder) Offset(offset int64) *Builder {
	b.sqlOffset = offset
	return b
}

func (b *Builder) OrderBy(columnName string, direction string) *Builder {
	if strings.EqualFold(direction, "desc") || strings.EqualFold(direction, "descending") {
		direction = "DESC"
	} else {
		direction = "ASC"
	}

	b.sqlOrderBy = append(b.sqlOrderBy, OrderBy{
		Column:    columnName,
		Direction: direction,
	})

	return b
}

func (b *Builder) Where(where Where) *Builder {
	b.sqlWhere = append(b.sqlWhere, where)
	return b
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

func (b *Builder) whereToSqlSingle(column string, operator string, value string) string {
	if operator == "==" || operator == "===" {
		operator = "="
	}
	if operator == "!=" || operator == "!==" {
		operator = "<>"
	}
	columnQuoted := b.quoteColumn(column)
	valueQuoted := b.quoteValue(value)

	sql := ""
	if b.Dialect == DIALECT_MYSQL {
		if value == "NULL" && operator == "=" {
			sql = columnQuoted + " IS NULL"
		} else if value == "NULL" && operator == "<>" {
			sql = columnQuoted + " IS NOT NULL"
		} else {
			sql = columnQuoted + " " + operator + " " + valueQuoted
		}
	}
	if b.Dialect == DIALECT_POSTGRES {
		if value == "NULL" && operator == "=" {
			sql = columnQuoted + " IS NULL"
		} else if value == "NULL" && operator == "<>" {
			sql = columnQuoted + " IS NOT NULL"
		} else {
			sql = columnQuoted + " " + operator + " " + valueQuoted
		}
	}
	if b.Dialect == DIALECT_SQLITE {
		if value == "NULL" && operator == "=" {
			sql = columnQuoted + " IS NULL"
		} else if value == "NULL" && operator == "<>" {
			sql = columnQuoted + " IS NOT NULL"
		} else {
			sql = columnQuoted + " " + operator + " " + valueQuoted
		}
	}
	return sql
}

/**
 * Converts wheres to SQL
 * @param array $wheres
 * @return string
 */
func (b *Builder) whereToSql(wheres []Where) string {
	sql := []string{}
	for _, where := range wheres {
		if where.Raw != "" {
			sql = append(sql, where.Raw)
			continue
		}

		if where.Type == "" {
			where.Type = "AND"
		}

		if where.Column != "" {
			sqlSingle := b.whereToSqlSingle(where.Column, where.Operator, where.Value)

			if len(sql) > 0 {
				sql = append(sql, where.Type+" "+sqlSingle)
			} else {
				sql = append(sql, sqlSingle)
			}

		}
		// } else {
		// 	$_sql = array();
		// 	$all = $where['WHERE'];
		// 	for ($k = 0; $k < count($all); $k++) {
		// 		$w = $all[$k];
		// 		$sqlSingle = $this->whereToSqlSingle($w['COLUMN'], $w['OPERATOR'], $w['VALUE']);
		// 		if ($k == 0) {
		// 			$_sql[] = $sqlSingle;
		// 		} else {
		// 			$_sql[] = $w['TYPE'] . " " . $sqlSingle;
		// 		}
		// 	}
		// 	$_sql = (count($_sql) > 0) ? " (" . implode(" ", $_sql) . ")" : "";

		// 	if ($i == 0) {
		// 		$sql[] = $_sql;
		// 	} else {
		// 		$sql[] = $where['TYPE'] . " " . $_sql;
		// 	}
		// }
	}

	if len(sql) > 0 {
		return " WHERE " + strings.Join(sql, " ")
	}

	return ""
}

//     private function groupby_to_sql($groupbys)
//     {
//         $sql = array();
//         // MySQL
//         if ($this->database_type == 'mysql') {
//             foreach ($groupbys as $groupby) {
//                 $sql[] = "`" . $groupby['COLUMN'] . "`";
//             }
//             return (count($sql) > 0) ? " GROUP BY " . implode(", ", $sql) : "";
//         }
//         // SQLite
//         if ($this->database_type == 'sqlite' or $this->database_type == 'sqlitedb') {
//             foreach ($groupbys as $groupby) {
//                 $sql[] = "" . $groupby['COLUMN'];
//             }
//             return (count($sql) > 0) ? " GROUP BY " . implode(", ", $sql) : "";
//         }
//     }

func (b *Builder) orderByToSql(orderBys []OrderBy) string {
	sql := []string{}

	if b.Dialect == DIALECT_MYSQL {
		for _, orderBy := range orderBys {
			sql = append(sql, b.quoteColumn(orderBy.Column)+" "+orderBy.Direction)
		}
	}

	if b.Dialect == DIALECT_POSTGRES {
		for _, orderBy := range orderBys {
			sql = append(sql, b.quoteColumn(orderBy.Column)+" "+orderBy.Direction)
		}
	}

	if b.Dialect == DIALECT_SQLITE {
		for _, orderBy := range orderBys {
			sql = append(sql, b.quoteColumn(orderBy.Column)+" "+orderBy.Direction)
		}
	}

	if len(sql) > 0 {
		return ` ORDER BY ` + strings.Join(sql, `,`)
	}

	return ""
}

func (b *Builder) quoteColumn(columnName string) string {
	columnSplit := strings.Split(columnName, ".")
	columnQuoted := []string{}

	for _, columnPart := range columnSplit {
		if b.Dialect == DIALECT_MYSQL {
			columnPart = "`" + columnPart + "`"
		}

		if b.Dialect == DIALECT_POSTGRES {
			columnPart = `"` + columnPart + `"`
		}

		if b.Dialect == DIALECT_SQLITE {
			columnPart = `"` + columnPart + `"`
		}

		columnQuoted = append(columnQuoted, columnPart)
	}

	return strings.Join(columnQuoted, ".")
}

func (b *Builder) quoteTable(tableName string) string {
	tableSplit := strings.Split(tableName, ".")
	tableQuoted := []string{}

	for _, tablePart := range tableSplit {
		if b.Dialect == DIALECT_MYSQL {
			tablePart = "`" + tablePart + "`"
		}

		if b.Dialect == DIALECT_POSTGRES {
			tablePart = `"` + tablePart + `"`
		}

		if b.Dialect == DIALECT_SQLITE {
			tablePart = `"` + tablePart + `"`
		}

		tableQuoted = append(tableQuoted, tablePart)
	}

	return strings.Join(tableQuoted, ".")
}

func (b *Builder) quoteValue(value string) string {
	if b.Dialect == DIALECT_MYSQL {
		value = `"` + value + `"`
	}

	if b.Dialect == DIALECT_POSTGRES {
		value = `"` + value + `"`
	}

	if b.Dialect == DIALECT_SQLITE {
		value = `'` + value + `'`
	}

	return value
}
