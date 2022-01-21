package orm

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[UPDATE] = _update
	generators[ORDERBY] = _orderby
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

func GetBindVars(num int) string {
	if num < 1 {
		return ""
	}
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ",")
}

func _insert(values ...interface{}) (string, []interface{}) {
	//values = tablename,(field1,field2...)
	tableName := values[0]
	fieldSql := strings.Join(values[1].([]string), ",")
	mainSql := fmt.Sprintf("insert into %s(%s)", tableName, fieldSql)
	return mainSql, []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	//values = (field1, field2)
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("values")
	for i, value := range values {
		sql.WriteString("(")
		v := value.([]interface{})
		bindVars := GetBindVars(len(v))
		sql.WriteString(bindVars)
		sql.WriteString(")")
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _select(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	fields := values[1].([]string)
	return fmt.Sprintf("select %s from %s", strings.Join(fields, ","), tableName), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	//limit := values[0].(int)
	return "limit ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("where %s", desc), vars
}

func _orderby(values ...interface{}) (string, []interface{}) {

	return fmt.Sprint("order by ", values[0].(string)), []interface{}{}
}

func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vs []interface{}
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vs = append(vs, v)
	}
	sql := fmt.Sprintf("update %s set %s", tableName, strings.Join(keys, ","))
	return sql, vs
}

func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("delete from %s", values[0].(string)), []interface{}{}
}

func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}
