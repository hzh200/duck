package dialects

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type sqlite3 struct{}

var _ Dialect = &sqlite3{}

func init() {
	RegisterDialect("sqlite3", &sqlite3{})
}

// sqlite_master
type SQLiteTableInfo struct {
	Typ      string
	Name     string
	TblName  string
	Rootpage int
	Sql      string
}

// The dialect doesn't support complex64, complex128, channel, uintptr, array, 
// map not using string as key or has value of multidimensional slice, slice which has more then two dimensions and most struct types.

func (s *sqlite3) DataTypeMapping(fieldType reflect.Type) string {
	switch fieldType.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return "INTEGER"
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return "UNSIGNED BIG INT"
	case reflect.Float32, reflect.Float64:
		return "REAL"
	case reflect.Bool:
		return "INTEGER"
	case reflect.String:
		return "TEXT"
	case reflect.Interface:
		return "TEXT"
	case reflect.Slice:
		return "TEXT"
	case reflect.Map:
		if fieldType.Key().Kind() != reflect.String {
			break
		}
		return "TEXT"
	case reflect.Struct:
		// if _, ok := fieldValue.Interface().(time.Time); ok {
		if fieldType == reflect.TypeOf(time.Time{}) {
			return "DATETIME"
		}
	}

	msg, _ := fmt.Printf("DataTypeMapping: data type %s not supported", fieldType)
	panic(msg)
}

func (s *sqlite3) DataFormatting(field interface{}) interface{} {
	switch reflect.ValueOf(field).Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return field
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return field
	case reflect.Float32, reflect.Float64:
		return field
	case reflect.Bool:
		return field
	case reflect.String:
		return fmt.Sprintf("'%s'", field)
	case reflect.Interface:
		res, _ := json.Marshal(field)
		return fmt.Sprintf("'%s'", res)
	case reflect.Slice:
		res, _ := json.Marshal(field)
		return fmt.Sprintf("'%s'", res)
	case reflect.Map:
		if reflect.TypeOf(field).Key().Kind() != reflect.String {
			break
		}
		res, _ := json.Marshal(field)
		return fmt.Sprintf("'%s'", res)
	case reflect.Struct:
		if _, ok := field.(time.Time); ok {
			return fmt.Sprintf("DateTime('%v')", field.(time.Time).Format(time.DateTime))
		}
	}

	msg, _ := fmt.Printf("DataFormatting: data type %s not supported", reflect.TypeOf(field))
	panic(msg)
}

func (s *sqlite3) DataDeformatting(value interface{}, fieldType reflect.Type) interface{} {
	switch fieldType.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return value
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return value
	case reflect.Float32, reflect.Float64:
		return value
	case reflect.Bool:
		return value
	case reflect.String:
		return value
	case reflect.Struct:
		if fieldType == reflect.TypeOf(time.Time{}) {
			return value
		}
	case reflect.Interface:
		var res interface{}
		json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
		return res
	}
	if fieldType.Kind() == reflect.Slice {
		eleType := fieldType.Elem()
		switch eleType.Kind() {
		case reflect.Int8:
			res := make([]int8, 0)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Int16:
			res := make([]int16, 0)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Int32:
			res := make([]int32, 0)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Int64:
			res := make([]int64, 0)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Int:
			res := make([]int, 0)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Uint8:
			res := make([]uint8, 0)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Uint16:
			res := make([]uint16, 0)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Uint32:
			res := make([]uint32, 0)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Uint64:
			res := make([]uint64, 0)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Uint:
			res := make([]uint, 0)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Float32:
			res := make([]float32, 0)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Float64:
			res := make([]float64, 0)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Bool:
			res := make([]bool, 0)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.String:
			res := make([]string, 0)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Slice:
			eleType = eleType.Elem()
			switch eleType.Kind() {
			case reflect.Int8:
				res := make([][]int8, 0)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Int16:
				res := make([][]int16, 0)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Int32:
				res := make([][]int32, 0)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Int64:
				res := make([][]int64, 0)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Int:
				res := make([][]int, 0)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Uint8:
				res := make([][]uint8, 0)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Uint16:
				res := make([][]uint16, 0)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Uint32:
				res := make([][]uint32, 0)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Uint64:
				res := make([][]uint64, 0)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Uint:
				res := make([][]uint, 0)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Float32:
				res := make([][]float32, 0)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Float64:
				res := make([][]float64, 0)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Bool:
				res := make([][]bool, 0)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.String:
				res := make([][]string, 0)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			}
		}
	}

	if fieldType.Kind() == reflect.Map {
		eleType := fieldType.Elem()
		switch eleType.Kind() {
		case reflect.Int8:
			res := make(map[string]int8)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Int16:
			res := make(map[string]int16)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Int32:
			res := make(map[string]int32)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Int64:
			res := make(map[string]int64)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Int:
			res := make(map[string]int)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Uint8:
			res := make(map[string]uint8)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Uint16:
			res := make(map[string]uint16)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Uint32:
			res := make(map[string]uint32)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Uint64:
			res := make(map[string]uint64)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Uint:
			res := make(map[string]uint)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Float32:
			res := make(map[string]float32)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Float64:
			res := make(map[string]float64)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Bool:
			res := make(map[string]bool)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.String:
			res := make(map[string]string)
			json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
			return res
		case reflect.Slice:
			eleType = eleType.Elem()
			switch eleType.Kind() {
			case reflect.Int8:
				res := make(map[string][]int8)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Int16:
				res := make(map[string][]int16)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Int32:
				res := make(map[string][]int32)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Int64:
				res := make(map[string][]int64)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Int:
				res := make(map[string][]int)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Uint8:
				res := make(map[string][]uint8)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Uint16:
				res := make(map[string][]uint16)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Uint32:
				res := make(map[string][]uint32)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Uint64:
				res := make(map[string][]uint64)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Uint:
				res := make(map[string][]uint)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Float32:
				res := make(map[string][]float32)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Float64:
				res := make(map[string][]float64)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.Bool:
				res := make(map[string][]bool)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			case reflect.String:
				res := make(map[string][]string)
				json.Unmarshal([]byte(reflect.ValueOf(value).Interface().(string)), &res)
				return res
			}
		}
	}

	msg, _ := fmt.Printf("DataDeformatting: data type %s not supported", fieldType)
	panic(msg)
}
