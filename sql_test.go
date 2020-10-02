package sql

import	"testing"


func TestCreation(t *testing.T) {
	sql := NewSqlite().Table("user").Select()
	if sql != "SELECT * FROM 'user';" {
		t.Fatalf(sql)
	}
}
