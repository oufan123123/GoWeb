package orm

import "reflect"

var dialectMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(value reflect.Value) string
	ExistTableSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectMap[name]
	return
}
