package orm

import (
	"fmt"
	"reflect"
	"time"
)

type Mysql struct{}

func init() {
	RegisterDialect("mysql", &Mysql{})
}

func (m *Mysql) DataTypeOf(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Bool:
		return "tinyint"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "varchar(255)"
	case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Uint32, reflect.Uintptr,
		reflect.Uint, reflect.Uint16, reflect.Uint8:
		return "int"
	case reflect.Struct:
		if _, ok := value.Interface().(time.Time); ok {
			return "date"
		}
	}
	panic(fmt.Sprintf("not found data type of mysql %s (%s)", value.Type().Name(), value.Kind()))
}

func (m *Mysql) ExistTableSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "select count(*) from information_schema.TABLES where TABLE_NAME = ?", args
}
