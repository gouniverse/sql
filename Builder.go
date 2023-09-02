package sql

import (
	"sort"
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

type GroupBy struct {
	Column string
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
	sqlGroupBy   []GroupBy
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
	if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
		sql = "DELETE FROM " + b.quoteTable(b.sqlTableName) + where + orderBy + limit + offset + ";"
	}
	return sql
}

// Drop deletes a table
func (b *Builder) Drop() string {
	sql := ""
	if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
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

func (b *Builder) GroupBy(groupBy GroupBy) *Builder {
	b.sqlGroupBy = append(b.sqlGroupBy, groupBy)
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

/** The <b>select</b> method selects rows from a table, based on criteria.
 * <code>
 * // Selects all the rows from the table
 * $db->table("USERS")->select();
 *
 * // Selects the rows where the column NAME is different from Peter, in descending order
 * $db->table("USERS")
 *     ->where("NAME","!=","Peter")
 *     ->orderby("NAME","desc")
 *     ->select();
 * </code>
 * @return mixed rows as associative array, false on error
 * @access public
 */
func (b *Builder) Select(columns []string) string {
	if b.sqlTableName == "" {
		panic("In method Delete() no table specified to delete from!")
	}

	join := "" // TODO

	groupBy := ""
	if len(b.sqlGroupBy) > 0 {
		groupBy = b.groupByToSql(b.sqlGroupBy)
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

	columnsStr := "*"

	if len(columns) > 0 {
		for index, column := range columns {
			columns[index] = b.quoteColumn(column)
		}
		columnsStr = strings.Join(columns, ", ")
	}

	sql := ""

	if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
		sql = "SELECT " + columnsStr + " FROM " + b.quoteTable(b.sqlTableName) + join + where + groupBy + orderBy + limit + offset + ";"
	}

	return sql
}

/**
 * The <b>update</b> method updates the values of a row in a table.
 * <code>
 * $updated_user = array("USER_MANE"=>"Mike");
 * $database->table("USERS")->where("USER_NAME","==","Peter")->update($updated_user);
 * </code>
 * @param Array an associative array, where keys are the column names of the table
 * @return int 0 or 1, on success, false, otherwise
 * @access public
 */
func (b *Builder) Insert(columnValuesMap map[string]string) string {
	if b.sqlTableName == "" {
		panic("In method Insert() no table specified to insert in!")
	}

	limit := ""
	if b.sqlLimit > 0 {
		limit = " LIMIT " + strconv.FormatInt(b.sqlLimit, 10)
	}

	offset := ""
	if b.sqlOffset > 0 {
		offset = " OFFSET " + strconv.FormatInt(b.sqlOffset, 10)
	}

	columnNames := []string{}
	columnValues := []string{}

	// Order keys
	keys := make([]string, 0, len(columnValuesMap))
	for k := range columnValuesMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, columnName := range keys {
		columnValue := columnValuesMap[columnName]
		columnNames = append(columnNames, b.quoteColumn(columnName))
		columnValues = append(columnValues, b.quoteValue(columnValue))
	}

	return "INSERT INTO " + b.quoteTable(b.sqlTableName) + " (" + strings.Join(columnNames, ", ") + ") VALUES (" + strings.Join(columnValues, ", ") + ")" + limit + offset + ";"
}

/**
 * The <b>update</b> method updates the values of a row in a table.
 * <code>
 * $updated_user = array("USER_MANE"=>"Mike");
 * $database->table("USERS")->where("USER_NAME","==","Peter")->update($updated_user);
 * </code>
 * @param Array an associative array, where keys are the column names of the table
 * @return int 0 or 1, on success, false, otherwise
 * @access public
 */
func (b *Builder) Update(columnValues map[string]string) string {
	if b.sqlTableName == "" {
		panic("In method Delete() no table specified to delete from!")
	}

	join := "" // TODO

	groupBy := ""
	if len(b.sqlGroupBy) > 0 {
		groupBy = b.groupByToSql(b.sqlGroupBy)
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

	// Order keys
	keys := make([]string, 0, len(columnValues))
	for k := range columnValues {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	updateSql := []string{}
	for _, columnName := range keys {
		columnValue := columnValues[columnName]
		updateSql = append(updateSql, b.quoteColumn(columnName)+"="+b.quoteValue(columnValue))
	}

	return "UPDATE " + b.quoteTable(b.sqlTableName) + " SET " + strings.Join(updateSql, ", ") + join + where + groupBy + orderBy + limit + offset + ";"
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

func (b *Builder) groupByToSql(groupBys []GroupBy) string {
	sql := []string{}
	for _, groupBy := range groupBys {
		sql = append(sql, b.quoteColumn(groupBy.Column))
	}

	if len(sql) > 0 {
		return " GROUP BY " + strings.Join(sql, ",")
	}

	return ""
}

// /**
//      * Joins tables to SQL.
//      * @return String the join SQL string
//      * @access private
//      */
// 	 private function join_to_sql($join, $table_name)
// 	 {
// 		 $sql = '';
// 		 // MySQL
// 		 if ($this->database_type == 'mysql') {
// 			 foreach ($join as $what) {
// 				 $type = $what[3] ?? '';
// 				 $alias = $what[4] ?? '';
// 				 $sql .= ' ' . $type . ' JOIN `' . $what[0] . '`';
// 				 if ($alias != "") {
// 					 $sql .= ' AS ' . $alias . '';
// 					 $what[0] = $alias;
// 				 }
// 				 if ($what[1] == $what[2]) {
// 					 $sql .= ' USING (`' . $what[1] . '`)';
// 				 } else {
// 					 $sql .= ' ON ' . $table_name . '.' . $what[1] . '=' . $what[0] . '.' . $what[2];
// 				 }
// 			 }
// 		 }
// 		 // SQLite
// 		 if ($this->database_type == 'sqlite' or $this->database_type == 'sqlitedb') {
// 			 foreach ($join as $what) {
// 				 $type = $what[3] ?? '';
// 				 $alias = $what[4] ?? '';
// 				 $sql .= " $type JOIN '" . $what[0] . "'";
// 				 if ($alias != "") {
// 					 $sql .= " AS '$alias'";
// 					 $what[0] = $alias;
// 				 }
// 				 $sql .= ' ON ' . $table_name . '.' . $what[1] . '=' . $what[0] . '.' . $what[2];
// 			 }
// 		 }

// 		 return $sql;
// 	 }

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
		value = `"` + b.escapeMysql(value) + `"`
	}

	if b.Dialect == DIALECT_POSTGRES {
		value = `"` + b.escapePostgres(value) + `"`
	}

	if b.Dialect == DIALECT_SQLITE {
		value = `'` + b.escapeSqlite(value) + `'`
	}

	return value
}

func (b *Builder) escapeMysql(value string) string {
	// escapeRegexp       = regexp.MustCompile(`[\0\t\x1a\n\r\"\'\\]`)
	// characterEscapeMap = map[string]string{
	// 	"\\0":  `\\0`,  //ASCII NULL
	// 	"\b":   `\\b`,  //backspace
	// 	"\t":   `\\t`,  //tab
	// 	"\x1a": `\\Z`,  //ASCII 26 (Control+Z);
	// 	"\n":   `\\n`,  //newline character
	// 	"\r":   `\\r`,  //return character
	// 	"\"":   `\\"`,  //quote (")
	// 	"'":    `\'`,   //quote (')
	// 	"\\":   `\\\\`, //backslash (\)
	// 	// "\\%":  `\\%`,  //% character
	// 	// "\\_":  `\\_`,  //_ character
	// }
	// return escapeRegexp.ReplaceAllStringFunc(val, func(s string) string {

	// 	mVal, ok := characterEscapeMap[s]
	// 	if ok {
	// 		return mVal
	// 	}
	// 	return s
	// })

	escapedStr := strings.ReplaceAll(value, `"`, `""`)
	return escapedStr
}

func (b *Builder) escapePostgres(value string) string {
	escapedStr := strings.ReplaceAll(value, "'", "''")
	return escapedStr
}

func (b *Builder) escapeSqlite(value string) string {
	escapedStr := strings.ReplaceAll(value, "'", "''")
	return escapedStr
}

/**
 * The <b>tables</b> method returns the names of all the tables, that
 * exist in the database.
 * <code>
 * foreach($database->tables() as $table){
 *     echo $table;
 * }
 * </code>
 * @param String the name of the table
 * @return array the names of the tables
 * @access public
 */
//  func (b *Builder) Tables(value string)
//  {
// 	 $tables = array();

// 	 if ($this->database_type == 'mysql') {
// 		 //$sql = "SHOW TABLES";
// 		 $sql = "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_TYPE='BASE TABLE' AND TABLE_SCHEMA='" . $this->database_name . "'";
// 		 $result = $this->executeQuery($sql);
// 		 if ($result === false)
// 			 return false;
// 		 foreach ($result as $row) {
// 			 $tables[] = $row['TABLE_NAME'];
// 		 }
// 		 return $tables;
// 	 }

// 	 if ($this->database_type == 'sqlite' or $this->database_type == 'sqlitedb') {
// 		 $sql = "SELECT * FROM 'SQLITE_MASTER' WHERE type='table' ORDER BY NAME ASC";
// 		 $result = $this->executeQuery($sql);
// 		 if ($result === false) {
// 			 return false;
// 		 }
// 		 foreach ($result as $row) {
// 			 $tables[] = $row['name'];
// 		 }
// 		 return $tables;
// 	 }
// 	 return false;
//  }
