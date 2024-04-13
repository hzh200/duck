package persistence

import "reflect"

type Dialect interface {
	DataTypeMapping(dbType reflect.Value)
}
