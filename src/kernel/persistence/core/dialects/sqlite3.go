package dialects

import (
	"reflect"
	"strings"
)

type sqlite3 struct {}

var _ Dialect = &sqlite3{}

func init() {
	RegisterDialect("sqlite3", &sqlite3{})
}

// sqlite_master
type SQLiteTableInfo struct {
	Typ string
	Name string
	TblName string
	Rootpage int
	Sql string
}

func (s *sqlite3) DataTypeMapping(fieldType reflect.Type) string {
	switch fieldType.Kind() {
		case reflect.String:
			return "TEXT"
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			return "INTEGER"
		case reflect.Float32, reflect.Float64:
			return "REAL"
		case reflect.Array, reflect.Slice:
			return "TEXT"
	}
	panic("cannot map struct type to model type")
}

func (s *sqlite3) DataFormatting(value interface{}) interface{} {
	fieldType := reflect.TypeOf(value)
	switch fieldType.Kind() {
		case reflect.Array, reflect.Slice:
			return strings.Join(value.([]string), ",")
	}
	return value
}

func (s *sqlite3) DataDeformatting(value interface{}) interface{} {
	fieldType := reflect.TypeOf(value)
	switch fieldType.Kind() {
		case reflect.Array, reflect.Slice:
			return strings.Join(value.([]string), ",")
	}
	return value
}
