package core

import (
	"duck/kernel/persistence/core/dialects"
	"reflect"
	"strings"
)

type Field struct {
	MemberName string
	MemberType reflect.Type
	FieldName string
	FieldType string
	Constraints string
	AutoIncrement bool
}

type Schema struct {
	ModelType reflect.Type
	StructName string
	TableName string
	Fields []Field	
	PrimaryKeyField Field
} 

func Parse(modelType reflect.Type, dialect dialects.Dialect, system bool) *Schema {
	schema := Schema{}
	schema.ModelType = modelType
	schema.StructName = modelType.Name()
	schema.TableName = toLowerCase(schema.StructName)
	schema.Fields = make([]Field, 0)
	for i := 0; i < modelType.NumField(); i++ {
		modelField := modelType.Field(i)
		// Embedding and non-exported fields.
		if modelField.Anonymous || !modelField.IsExported() {
			continue
		}
		field := Field{}
		field.MemberName = modelField.Name
		field.FieldName = toLowerCase(field.MemberName)
		field.MemberType = modelField.Type
		field.FieldType = dialect.DataTypeMapping(modelField.Type)
		field.Constraints = modelField.Tag.Get("constraints");
		field.AutoIncrement = strings.Contains(field.Constraints, "AutoIncrement")

		if strings.Contains(field.Constraints, "PrimaryKey") {
			if schema.PrimaryKeyField.MemberName != "" {
				panic("model should only contain one primary key")
			}
			schema.PrimaryKeyField = field
		}

		schema.Fields = append(schema.Fields, field)
	}

	if !system && schema.PrimaryKeyField.MemberName == ""  {
		panic("model must contain a primary key")
	}
	
	return &schema
}

func toLowerCase(name string) string {
	lowerCasedNameBuilder := strings.Builder{}
	for i, c := range []byte(name) {
		if c >= 65 && c <= 90 {
			c = c + 32
			if i != 0 {
				lowerCasedNameBuilder.WriteRune('_')
			}
		}
		lowerCasedNameBuilder.WriteByte(c)
	}
	return lowerCasedNameBuilder.String()
}
