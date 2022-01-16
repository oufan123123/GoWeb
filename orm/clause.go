package orm

import (
	"strings"
)

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

type Clause struct {
	Sql     map[Type]string
	SqlVars map[Type][]interface{}
}

func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.Sql == nil {
		c.Sql = make(map[Type]string)
		c.SqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.Sql[name] = sql
	c.SqlVars[name] = vars
}

func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var sqlVars []interface{}
	for _, order := range orders {
		if sql, ok := c.Sql[order]; ok {
			sqls = append(sqls, sql)
			sqlVars = append(sqlVars, c.SqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), sqlVars
}
