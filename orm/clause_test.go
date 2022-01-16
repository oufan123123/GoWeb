package orm

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSelect(t *testing.T) {
	var c Clause
	c.Set(LIMIT, 2)
	c.Set(SELECT, "user_name", "user")
	c.Set(WHERE, "is_locked = ?", true)
	c.Set(ORDERBY, "age desc")
	sql, vars := c.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)
	fmt.Println(sql)
	if sql != "select user_name from user where is_locked = ? order by age desc limit ?" {
		t.Fatal("error creating select sql")
		//Info(sql)
	}
	if !reflect.DeepEqual(vars, []interface{}{true, 1}) {
		t.Fatal("error creating select sql of vars")
	}
}

func TestClause(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		TestSelect(t)
	})
}
