package sql

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/emirpasic/gods/stacks/arraystack"

	"github.com/emirpasic/gods/maps/hashmap"
)

// SQLBuilder represents an SQL builder
type SQLBuilder struct {
	Type      string
	Cmd       string
	tableName string
	//SQL       hashmap.Map
	where *arraystack.Stack
	// TagContent    string
	// TagAttributes map[string]string
	// TagChildren   []*Tag
}

// NewSqlite represents a BUTTON tag
func NewSqlite() *SQLBuilder {
	sql := &SQLBuilder{
		Type:  "sqlite",
		where: arraystack.New(),
	}
	return sql
}

// Addslashes addslashes()
func addslashes(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '\'', '"', '\\':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// Addslashes addslashes()
// func addslashes(str string) string {
// 	var buf bytes.Buffer
// 	for _, char := range str {
// 		switch char {
// 		case '\'', '"', '\\':
// 			buf.WriteRune('\\')
// 		}
// 		buf.WriteRune(char)
// 	}
// 	return buf.String()
// }

func inArrayString(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// Insert returns the SQL query
func (sqlBuilder *SQLBuilder) Insert(columnValueMap map[string]string) string {
	sql := ""
	numOfColumns := len(columnValueMap)
	columnNames := make([]string, numOfColumns)
	columnValues := make([]string, numOfColumns)

	i := 0
	for columnName, columnValue := range columnValueMap {
		fmt.Println("Key:", columnName, "Value:", columnValue)
		columnValue = sqlBuilder.escapeSqlite(columnValue)
		columnNames[i] = columnName
		columnValues[i] = columnValue
		i++
	}

	columnNamesStr := "'" + strings.Join(columnNames, "','") + "'"
	columnValuesStr := "'" + strings.Join(columnValues, "','") + "'"
	//tableNameInterface, _ := sqlBuilder.SQL.Get("TableName")
	//tableName := fmt.Sprintf("%v", tableNameInterface)

	sql = "INSERT INTO '" + sqlBuilder.tableName + "'(" + columnNamesStr + ") VALUES (" + columnValuesStr + ")"

	// if ($this->database_type == 'sqlite' or $this->database_type == 'sqlitedb') {
	//     foreach ($row_values as $key => $value) {
	//         $row_values[$key] = is_null($value) ? 'NULL' : $this->dbh->quote($value);
	//     }
	//     $values = implode(",", array_values($row_values));
	//     $fields = implode("','", array_keys($row_values));
	//     $sql = "INSERT INTO '" . $table_name . "'('" . $fields . "') VALUES (" . $values . ")";
	// }

	// $this->open();
	//     if (isset($this->sql["table"]) == false)
	//         trigger_error('ERROR: In class <b>' . get_class($this) . '</b> in method <b>insert($row_values)</b>: Not specified table to insert a row in!', E_USER_ERROR);
	//     if (is_array($row_values) == false)
	//         trigger_error('ERROR: In class <b>' . get_class($this) . '</b> in method <b>insert($row_values)</b>: Parameter <b>$row_values</b> MUST BE of type Array - <b style="color:red">' . gettype($row_values) . '</b> given!', E_USER_ERROR);
	//     $table_name = $this->sql["table"][0];

	//     if ($this->database_type == 'mysql') {
	//         foreach ($row_values as $key => $value) {
	//             $row_values[$key] = is_null($value) ? 'NULL' : $this->dbh->quote($value);
	//         }
	//         $values = implode(",", array_values($row_values));
	//         $fields = "`" . implode("`" . "," . "`", array_keys($row_values)) . "`";
	//         $sql = 'INSERT INTO `' . $table_name . '`(' . $fields . ') VALUES (' . $values . ')';
	//     }

	//     if ($this->database_type == 'sqlite' or $this->database_type == 'sqlitedb') {
	//         foreach ($row_values as $key => $value) {
	//             $row_values[$key] = is_null($value) ? 'NULL' : $this->dbh->quote($value);
	//         }
	//         $values = implode(",", array_values($row_values));
	//         $fields = implode("','", array_keys($row_values));
	//         $sql = "INSERT INTO '" . $table_name . "'('" . $fields . "') VALUES (" . $values . ")";
	//     }
	//     $this->sql = array(); // Emptying the SQL array
	//     $result = $this->executeNonQuery($sql);
	//     if ($result === false)
	//         return false;
	//     return $result;

	sqlBuilder.Empty()

	return sql
}

// Table returns the SQL query
func (sqlBuilder *SQLBuilder) Table(tableName string) *SQLBuilder {
	//sqlBuilder.SQL.Put(tableName, "TableName")
	sqlBuilder.tableName = tableName
	return sqlBuilder
}

// Select returns a SELECT query
func (sqlBuilder *SQLBuilder) Select() string {
	if sqlBuilder.Type == "sqlite" {
		return sqlBuilder.selectSqlite()
	}
	return ""
	// $sql = '';

	//     if (isset($this->sql["table"]) == false) {
	//         trigger_error('ERROR: In class <b>' . get_class($this) . '</b> in method <b>select()</b>: Not specified table to select from!', E_USER_ERROR);
	//     }

	//     $table_name = $this->sql["table"][0];
	//     $where = isset($this->sql["where"]) == false ? '' : $this->where_to_sql($this->sql["where"]);
	//     $orderby = isset($this->sql["orderby"]) == false ? '' : $this->orderby_to_sql($this->sql["orderby"]);
	//     $limit = (isset($this->sql["limit"]) == false) ? '' : " LIMIT " . $this->sql["limit"];
	//     $groupby = isset($this->sql["groupby"]) == false ? '' : $this->groupby_to_sql($this->sql["groupby"]);
	//     $join = isset($this->sql["join"]) == false ? '' : $this->join_to_sql($this->sql["join"], $table_name);
	//     if (is_array($columns)) {
	//         if (count($columns) > 0) {
	//             $columns = implode(',', $columns);
	//         }
	//     }

	//     if ($this->database_type == 'mysql') {
	//         $sql = 'SELECT ' . $columns . ' FROM `' . $table_name . '`' . $join . $where . $groupby . $orderby . $limit . ';';
	//     }

	//     if ($this->database_type == 'sqlite' or $this->database_type == 'sqlitedb') {
	//         $sql = "SELECT " . $columns . " FROM '" . $table_name . "'" . $join . $where . $groupby . $orderby . $limit . ";";
	//     }

	//     if ($this->sqlOutput) {
	//         $this->sqlOutput = false; // Disable for future queries
	//         return $sql;
	//     }

	//     $this->sql = array(); // Emptying the SQL array
	//     $result = $this->executeQuery($sql);
	//     if ($result === false) {
	//         return false;
	//     }
	//     return $result;
}

func (sqlBuilder *SQLBuilder) selectSqlite() string {
	where := sqlBuilder.whereToSQL()
	columns := "*"
	sql := "SELECT "
	sql += columns
	sql += " FROM '" + sqlBuilder.tableName + "'"
	// sql += join
	sql += where
	// sql += groupby
	// sql += orderby
	// sql += limit
	sql += ";"

	return sql
}

// Where adds a WHERE clause
func (sqlBuilder *SQLBuilder) Where(columnName string, comparisonOperator string, value string) *SQLBuilder {
	entry := hashmap.New()
	entry.Put("column", columnName)
	entry.Put("operator", comparisonOperator)
	entry.Put("value", value)
	entry.Put("type", "AND")

	sqlBuilder.where.Push(entry)

	return sqlBuilder
}

// Empty empries the query
func (sqlBuilder *SQLBuilder) Empty() {
	sqlBuilder.tableName = ""
	sqlBuilder.where = arraystack.New()
}

// whereToSQL creates a WHERE clause
func (sqlBuilder *SQLBuilder) whereToSQL() string {
	if sqlBuilder.Type == "sqlite" {
		return sqlBuilder.whereToSQLSqlite()
	}
	return "1=1"
}

func (sqlBuilder *SQLBuilder) whereToSQLSqlite() string {
	sqls := make([]string, sqlBuilder.where.Size())
	it := sqlBuilder.where.Iterator()
	for it.Next() {
		index, value := it.Index(), it.Value()
		where := value.(*hashmap.Map)
		whereColumn, _ := where.Get("column")
		whereOperator, _ := where.Get("operator")
		whereValue, _ := where.Get("value")
		whereType, _ := where.Get("type")
		whereSQL := sqlBuilder.whereToSQLSqliteSingle(whereColumn.(string), whereOperator.(string), whereValue.(string))
		if index != 0 {
			whereSQL = whereType.(string) + " " + whereSQL
		}
		sqls[index] = whereSQL
	}
	sql := strings.Join(sqls, " ")
	if len(sql) > 0 {
		return " WHERE " + sql
	}
	return ""
}

func (sqlBuilder *SQLBuilder) whereToSQLSqliteSingle(columnName string, operator string, value string) string {
	sql := ""
	// $column = explode('.', $column);
	columnQuoted := "" + columnName + ""
	if operator == "==" || operator == "===" {
		operator = "="
	}
	if operator == "!=" || operator == "!==" {
		operator = "<>"
	}
	if value == "{{NULL}}" && operator == "=" {
		sql = columnQuoted + " IS NULL"
	} else if value == "{{NULL}}" && operator == "<>" {
		sql = columnQuoted + " IS NOT NULL"
	} else {
		valueQuoted := sqlBuilder.quoteSqlite(value)
		sql = columnQuoted + " " + operator + " " + valueQuoted
	}
	return sql
}

func (sqlBuilder *SQLBuilder) escapeSqlite(value string) string {
	escapedStr := strings.ReplaceAll(value, "'", "''")
	return escapedStr
}
func (sqlBuilder *SQLBuilder) quoteSqlite(value string) string {
	quotedStr := "'" + sqlBuilder.escapeSqlite(value) + "'"
	return quotedStr
}
