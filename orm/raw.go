package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Session struct {
	db       *sql.DB
	dialect  Dialect
	clause   Clause
	refTable *Schema
	sql      strings.Builder
	sqlVars  []interface{}
}

func NewSession(db *sql.DB, d Dialect) *Session {
	return &Session{
		db:      db,
		dialect: d,
	}
}

func (s *Session) getRefTable() *Schema {
	if s.refTable == nil {
		Error("model not set")
	}
	return s.refTable
}

func (s *Session) Model(obj interface{}) {

	if s.refTable == nil || reflect.TypeOf(obj) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = Parse(obj, s.dialect)
		return
	}
	Info("model reftable is null or table has been modeled")
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = Clause{}
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) CreateTable() error {
	refTable := s.getRefTable()
	var columns []string
	for _, field := range refTable.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	column := strings.Join(columns, ",")
	column = fmt.Sprintf("create table %s(%s);", refTable.Name, column)
	_, err := s.Raw(column).Exec()
	return err
}

func (s *Session) DropTable() error {
	refTable := s.getRefTable()
	_, err := s.Raw(fmt.Sprintf("drop table if exists %s", refTable.Name)).Exec()
	return err
}

func (s *Session) HasTable() error {
	refTable := s.getRefTable()
	sql, sqlVars := s.dialect.ExistTableSQL(refTable.Name)
	row := s.Raw(sql, sqlVars...).QueryRow()
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
		//Errorf("scan table error%s",refTable.Name)
	} else {
		if count > 0 {
			fmt.Printf("find table(%s) in database", refTable.Name)
		} else {
			fmt.Printf("not find table(%s) in database", refTable.Name)
		}
		return nil

	}

}

func (s *Session) Raw(sql string, sqlvars ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, sqlvars...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		Error(err)
	}
	return
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)

}

func (s *Session) Insert(objects ...interface{}) (int64, error) {
	if objects == nil || len(objects) < 1 {
		Info("insert obj is null")
		return 0, fmt.Errorf("insert obj is null line affected:%d", 0)
	}
	var res []interface{}
	for _, object := range objects {
		// why many times?
		s.Model(object)
		table := s.getRefTable()
		s.clause.Set(INSERT, table.Name, table.FieldNames)
		res = append(res, s.getRefTable().RecordValues(object))
	}
	s.clause.Set(VALUES, res...)
	sql, vars := s.clause.Build(INSERT, VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Find(object interface{}) error {
	//s.Model((*object)[0])
	obj := reflect.Indirect(reflect.ValueOf(object))
	elemType := obj.Type().Elem()
	s.Model(reflect.New(elemType).Elem().Interface())
	table := s.getRefTable()
	s.clause.Set(SELECT, table.Name, table.FieldNames)
	sql, _ := s.clause.Build(SELECT)
	results, err := s.Raw(sql).QueryRows()
	if err != nil {
		return err
	}

	for results.Next() {
		e := reflect.New(elemType).Elem()
		var values []interface{}
		for _, field := range table.Fields {
			values = append(values, e.FieldByName(field.Name).Addr().Interface())
		}
		if err := results.Scan(values...); err != nil {
			return err
		}
		obj.Set(reflect.Append(obj, e))
	}
	return results.Close()
}
