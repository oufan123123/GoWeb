package orm

import (
	"fmt"
	"go/ast"
	"reflect"
)

// reflect to a mysql table field
type Field struct {
	Name string
	Type string
	Tag  string
}

// reflect to a mysql table
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	FieldMap   map[string]*Field
}

func Parse(object interface{}, d Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(object)).Type()
	schema := &Schema{
		Model:    object,
		Name:     modelType.Name(),
		FieldMap: make(map[string]*Field),
	}
	x := modelType.NumField()
	fmt.Print(x)
	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i)
		if !f.Anonymous && ast.IsExported(f.Name) {
			field := &Field{
				Name: f.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(f.Type))),
			}
			if v, ok := f.Tag.Lookup("orm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, field.Name)
			schema.FieldMap[field.Name] = field
		}
	}
	return schema
}

func (schema *Schema) GetField(name string) *Field {
	return schema.FieldMap[name]
}

func (schema *Schema) RecordValues(object interface{}) []interface{} {

	//Parse(objects[0])
	var res []interface{}
	model := reflect.Indirect(reflect.ValueOf(object))
	for _, field := range schema.Fields {
		res = append(res, model.FieldByName(field.Name).Interface())
	}
	return res
}
