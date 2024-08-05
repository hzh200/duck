package dialects

import "reflect"

type Dialect interface {
	DataTypeMapping(fieldType reflect.Type) string
	DataFormatting(field interface{}) interface{}
	DataDeformatting(value interface{}, fieldType reflect.Type) interface{}
	AutoIncrement() string
	PrimaryKey() string
}

var dialects map[string]Dialect = map[string]Dialect{}

func RegisterDialect(name string, dialect Dialect) {
	dialects[name] = dialect
}

func GetDialect(name string) Dialect {
	if dialect, ok := dialects[name]; ok {
		return dialect
	} else {
		return nil
	}
}
